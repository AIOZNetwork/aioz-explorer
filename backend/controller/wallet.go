package controller

import (
	"aioz.io/go-aioz/types"
	"aioz.io/go-aioz/x_gob_explorer/config"
	"aioz.io/go-aioz/x_gob_explorer/context"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"aioz.io/go-aioz/x_gob_explorer/utils"
	"encoding/hex"
	cmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"strings"
)

func (ctrl *Controller) parseWallets(ctx context.Context, listWallets map[string]bool) error {
	coinDenom := config.GetConfig().GetString("coins.denom")
	stakeDenom := config.GetConfig().GetString("stake.denom")

	walletAccAddresses := make([]*entity.WalletAddress, 0)
	walletValidators := make([]*entity.Validator, 0)
	for a, v := range listWallets {
		if !v {
			continue
		}
		switch getWalletPrefix(a) {
		case types.Bech32PrefixValAddr:
			valAddress, err := cmtypes.ValAddressFromBech32(a)
			if err != nil {
				return err
			}
			val, found := stakingKeeper.GetValidator(ctx.CMCtx(), valAddress)
			if !found {
				continue
			}
			object := &entity.Validator{
				Address: a,
				Tokens:  val.Tokens.String(),
				Power:   val.GetConsensusPower(),
				Jailed:  val.Jailed,
				Status:  val.Status.String(),
			}
			walletValidators = append(walletValidators, object)
			
		default:
			accAddress, err := cmtypes.AccAddressFromBech32(a)
			if err != nil {
				return err
			}
			acc := accountKeeper.GetAccount(ctx.CMCtx(), accAddress)
			if acc == nil {
				// this address does not exists on blockchain yet,
				//just available temporary on GoB
				continue
			}
			pubkey := ""
			if acc.GetPubKey() != nil {
				pubkey = hex.EncodeToString(acc.GetPubKey().Bytes())
			}
			coinAIOZ := ""
			for _, c := range acc.GetCoins() {
				if c.Denom == coinDenom {
					coinAIOZ = c.Amount.String()
					break
				}
			}
			decAIOZ, _ := decimal.NewFromString(coinAIOZ)
			coinStake := ""
			for _, c := range acc.GetCoins() {
				if c.Denom == stakeDenom {
					coinStake = c.Amount.String()
					break
				}
			}
			decStake, _ := decimal.NewFromString(coinStake)
			w := &entity.WalletAddress{
				Address:       a,
				Coins:         acc.GetCoins().String(),
				AccountNumber: acc.GetAccountNumber(),
				Sequence:      acc.GetSequence(),
				PubKey:        pubkey,
				CoinAIOZ:      decAIOZ,
				CoinStake:     decStake,
			}
			walletAccAddresses = append(walletAccAddresses, w)		
		}
	}
	limit := 1500
	for i := 0; i < len(walletValidators); i += limit {
		batch := walletValidators[i:utils.Min2Int(i+limit, len(walletValidators))]
		if err := ctrl.validatorRepo.MultiRowsUpsert(ctx, batch); err != nil {
			return err
		}
	}

	for i := 0; i < len(walletAccAddresses); i += limit {
		batch := walletAccAddresses[i:utils.Min2Int(i+limit, len(walletAccAddresses))]
		if err := ctrl.walletRepo.MultiRowsUpsert(ctx, batch); err != nil {
			return err
		}
	}
	return nil
}

func getWalletPrefix(address string) string {
	tokens := strings.Split(address, "-")
	if len(tokens) <= 1 {
		return ""
	}
	return tokens[0] + "-"
}
