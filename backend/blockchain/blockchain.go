package blockchain

import (
	"aioz.io/go-aioz/x_gob_explorer/config"
	"aioz.io/go-aioz/x_gob_explorer/context"
	ctrl "aioz.io/go-aioz/x_gob_explorer/controller"
	"aioz.io/go-aioz/x_gob_explorer/email"
	"aioz.io/go-aioz/x_gob_explorer/log"
	"aioz.io/go-aioz/x_gob_explorer/ws"
	cmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/labstack/echo/v4"
	"github.com/tendermint/tendermint/types"
	"strings"
)

var (
	blockCh    = make(chan *types.Block, 1000000)
	ctxCreator func(h int64) cmtypes.Context
	controller ctrl.Controller
)

func Init(ctrl ctrl.Controller, cc func(height int64) cmtypes.Context) {

	e := echo.New()
	controller = ctrl
	ctxCreator = cc

	wshandler := ws.NewWS()
	e.GET("/ws/join", wshandler.ServeWebsocket)
	go wshandler.Run()
	go StartIndexing(blockCh, wshandler)

	wsserverurl := config.GetConfig().GetString("ws.url")
	e.Logger.Fatal(e.Start(wsserverurl))
}

func ListenNewBlock(block *types.Block) {
	blockCh <- block
}

func getHighestBlockInDB() (int64, error) {
	return controller.GetHighestBlockInDB()
}

func StartIndexing(blockCh chan *types.Block, wshandler *ws.WS) {
	log.InitLogger()
	logger := log.GoBLogger()
	b, err := getHighestBlockInDB()
	if err != nil {
		// handle if no block exists -> 404 not found
		// just set current block = 0
		if strings.Contains(err.Error(), "not found") {
			b = 0
		}
	}
	for block := range blockCh {
		if block.Height <= b {
			continue
		}
		cmContext := ctxCreator(block.Height)
		ctx := context.NewGoBContext(block.Height, logger, cmContext)

		objTxs, err := controller.ProcessingIndex(ctx, block)
		if err != nil {
			logger.Error(err)
			_ = email.SendMail(err.Error())
			panic(err)
		}

		if err := controller.ExecuteTxn(); err != nil {
			logger.Error(err)
			_ = email.SendMail(err.Error())
			panic(err)
		}

		// handle ws subscribe
		subscribeWallets(wshandler, objTxs)
		subscribeMessageSend(wshandler, objTxs)
		subscribeBlock(wshandler, &struct {
			Height  int64  `json:"height"`
			Time    int64  `json:"time"`
			ChainId string `json:"chain_id"`
		}{Height: block.Height, Time: block.Time.Unix(), ChainId: block.ChainID})
	}
}
