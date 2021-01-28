package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"net/http"
	"swagger-server/context"
	"swagger-server/domain"
)

type statUsecase struct {
	statRepo           domain.StatisticRepository
	msgsUsecase        domain.MsgsUsecase
	blockUsecase       domain.BlockUsecase
	transactionUsecase domain.TransactionUsecase
	walletUsecase      domain.WalletUsecase
	lcdUsecase         domain.LCDUsecase
	cdc                *codec.Codec
}

func NewStatUsecase(statRepo domain.StatisticRepository, mu domain.MsgsUsecase, bu domain.BlockUsecase,
	tu domain.TransactionUsecase, wu domain.WalletUsecase, lu domain.LCDUsecase, cdc *codec.Codec) domain.StatisticUsecase {
	return &statUsecase{
		statRepo:           statRepo,
		msgsUsecase:        mu,
		blockUsecase:       bu,
		transactionUsecase: tu,
		walletUsecase:      wu,
		lcdUsecase:         lu,
		cdc:                cdc,
	}
}

func (s *statUsecase) GetInflation(ctx context.Context) (string, error) {
	host := viper.GetString("lcdserver.host")
	port := viper.GetString("lcdserver.port")
	res, body, err := s.lcdUsecase.Request(host, port, "GET", "/minting/inflation", nil)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", errors.New("response status not ok")
	}
	result, err := s.lcdUsecase.ExtractResultFromResponse([]byte(body))
	if err != nil {
		return "", err
	}
	var inflation string
	if err := s.cdc.UnmarshalJSON(result, &inflation); err != nil {
		return "", err
	}
	dec, _ := decimal.NewFromString(inflation)
	dec = dec.Mul(decimal.NewFromInt(100))
	return fmt.Sprintf("%v%%", dec.StringFixed(2)), nil
}

func (s *statUsecase) GetTotalSupply(ctx context.Context) (string, error) {
	host := viper.GetString("lcdserver.host")
	port := viper.GetString("lcdserver.port")
	res, body, err := s.lcdUsecase.Request(host, port, "GET", "/supply/total", nil)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", errors.New("response status not ok")
	}
	result, err := s.lcdUsecase.ExtractResultFromResponse([]byte(body))
	if err != nil {
		return "", err
	}
	var ttsupply sdk.Coins
	if err := s.cdc.UnmarshalJSON(result, &ttsupply); err != nil {
		return "", err
	}
	b, _ := json.Marshal(ttsupply)
	return string(b), nil
}

func (s *statUsecase) GetStatistic(ctx context.Context) (*domain.Statistic, error) {
	volume24h, err := s.msgsUsecase.GetVolumeSend24h(ctx)
	if err != nil {
		return nil, err
	}
	numBlocks, err := s.blockUsecase.CountTotalBlocks(ctx)
	if err != nil {
		return nil, err
	}
	numTx, err := s.transactionUsecase.CountTotalTx(ctx)
	if err != nil {
		return nil, err
	}
	bestBlock, err := s.blockUsecase.GetBestBlock(ctx)
	if err != nil {
		return nil, err
	}
	topWallets, err := s.walletUsecase.GetTopWallet(ctx, 5)
	if err != nil {
		return nil, err
	}
	totalWallets, err := s.walletUsecase.CountTotalWallets(ctx)
	if err != nil {
		return nil, err
	}
	cirSupply, err := sdk.ParseCoins(bestBlock.CirculatingSupply)
	if err != nil {
		return nil, err
	}
	totalSupply, err := sdk.ParseCoins(bestBlock.TotalSupply)
	if err != nil {
		return nil, err
	}
	a, _ := json.Marshal(cirSupply)
	b, _ := json.Marshal(totalSupply)
	inflation := bestBlock.AnnualInflationRate.Mul(decimal.NewFromInt(100))
	inf := fmt.Sprintf("%v%%", inflation.String())
	obj := &domain.Statistic{
		TotalBlock:       numBlocks,
		TotalTransaction: numTx,
		TotalWallets:     totalWallets,
		Volume24h:        volume24h,
		//Volume24h:        "haha",
		AvgFee:           bestBlock.AvgFee,
		TopAccount:       topWallets,
		MarketCap:        "$9,350,238",
		BlockTime:        "5 seconds",
		Inflation:         inf,
		CirculatingSupply: string(a),
		TotalSupply:       string(b),
	}
	return obj, nil
}
