package controller

import (
	"aioz.io/go-aioz/simapp"
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/db"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"aioz.io/go-aioz/x_gob_explorer/domain/repository"
	"aioz.io/go-aioz/x_gob_explorer/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	cmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/types"
)

var (
	cdc                *codec.Codec
	txDecoder          cmtypes.TxDecoder
	accountKeeper      auth.AccountKeeper
	stakingKeeper      staking.Keeper
	supplyKeeper       supply.Keeper
	distributionKeeper distribution.Keeper
)

type Controller struct {
	db              db.Database
	blockRepo       repository.BlockRepo
	transactionRepo repository.TransactionRepo
	messageRepo     repository.MessageRepo
	txsRepo         repository.TxsRepo
	walletRepo      repository.WalletRepo
	validatorRepo   repository.ValidatorRepo
	delegatorRepo   repository.DelegatorRepo
	stakeRepo       repository.StakedRepo
	gs              simapp.GenesisState
}

func NewController(db db.Database,
	codec *codec.Codec, txDec cmtypes.TxDecoder, accKeeper auth.AccountKeeper, stkKeeper staking.Keeper,
	suppKeeper supply.Keeper, distrKeeper distribution.Keeper) Controller {
	cdc = codec
	txDecoder = txDec
	accountKeeper = accKeeper
	stakingKeeper = stkKeeper
	supplyKeeper = suppKeeper
	distributionKeeper = distrKeeper

	return Controller{
		db:              db,
		blockRepo:       repository.NewBlockRepo(db),
		transactionRepo: repository.NewTransactionRepo(db),
		messageRepo:     repository.NewMessageRepo(db),
		txsRepo:         repository.NewTxsRepo(db),
		walletRepo:      repository.NewWalletRepo(db),
		validatorRepo:   repository.NewValidatorRepo(db),
		delegatorRepo:   repository.NewDelegatorRepo(db),
		stakeRepo:       repository.NewStakedRepo(db),
	}
}

func (ctrl *Controller) SetSimApp(gs simapp.GenesisState) {
	ctrl.gs = gs
}

func (ctrl *Controller) ProcessingIndex(ctx context.Context, block *types.Block) ([]*entity.Txs, error) {
	// processing block-by-block
	objBlock := ctrl.preprocessBlockData(ctx, block)
	if err := ctrl.blockRepo.CreateBlock(ctx, objBlock); err != nil {
		return nil, err
	}

	// processing all txs in block
	txns, signer, err := ctrl.preprocessTransactionData(ctx, block)
	if err != nil {
		return nil, err
	}

	limit := 1000
	for i := 0; i < len(txns); i += limit {
		batch := txns[i:utils.Min2Int(i+limit, len(txns))]
		if err = ctrl.transactionRepo.MultiRowsInsert(ctx, batch); err != nil {
			return nil, err
		}
	}
	// processing all messages in txs in block
	msgWallets, objTxs, err := ctrl.preprocessMessageData(ctx, block)
	if err != nil {
		return nil, err
	}
	// processing each validator's reward pool
	err = ctrl.handleValidatorsPool(ctx)
	if err != nil {
		return nil, err
	}

	// all wallets related to block, but depends on type of messages received
	appendWalletsToMapWallet(msgWallets, signer)
	if err := ctrl.parseWallets(ctx, msgWallets); err != nil {
		return nil, err
	}

	return objTxs, nil
}

func (ctrl *Controller) ExecuteTxn() error {
	return ctrl.db.FlushAllTxns()
}

func (ctrl *Controller) GetHighestBlockInDB() (int64, error) {
	block, err := ctrl.blockRepo.GetPreviousBlockInfo()
	if err != nil {
		return -1, err
	}
	return block.Height, nil
}
