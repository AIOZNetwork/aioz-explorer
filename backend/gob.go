package x_gob_explorer

import (
	"aioz.io/go-aioz/simapp"
	"aioz.io/go-aioz/x_gob_explorer/blockchain"
	"aioz.io/go-aioz/x_gob_explorer/config"
	"aioz.io/go-aioz/x_gob_explorer/controller"
	gobdb "aioz.io/go-aioz/x_gob_explorer/domain/db"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"aioz.io/go-aioz/x_gob_explorer/email"
	"github.com/cosmos/cosmos-sdk/codec"
	cmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	crdb gobdb.Database
)

func init() {

	// load config and create connection
	config.LoadConfig("config/config.toml")

	host := config.GetConfig().GetString("database.host")
	port := config.GetConfig().GetString("database.port")
	user := config.GetConfig().GetString("database.user")
	pwd := config.GetConfig().GetString("database.passwd")
	dbname := config.GetConfig().GetString("database.dbname")
	sslmode := config.GetConfig().GetString("database.sslmode")
	sslrootcert := config.GetConfig().GetString("database.sslrootcert")
	sslkey := config.GetConfig().GetString("database.sslkey")
	sslcert := config.GetConfig().GetString("database.sslcert")

	crdb = gobdb.NewCockroachDB(user, pwd, host, port, dbname, sslmode, sslrootcert, sslkey, sslcert)
}

func prepareDB(client *gorm.DB) {
	// init tables
	start := time.Now()
	if err := client.AutoMigrate(
		&entity.Block{}, &entity.Transaction{}, &entity.Txs{}, &entity.MessageSend{},
		&entity.WalletAddress{}, &entity.Validator{}, &entity.Delegator{}, &entity.Stake{}, &entity.NodeInfo{}, &entity.PnTokenDevice{}); err != nil {
		_ = email.SendMail(err.Error())
		panic(err)
	}
	log.Printf("Auto migrate all tables takes: %v(ms)\n", time.Now().Sub(start).Milliseconds())
}

func InitGOB(ctxCreator func(height int64) cmtypes.Context, cdc *codec.Codec, txDec cmtypes.TxDecoder, accKeeper auth.AccountKeeper,
	stkKeeper staking.Keeper, suppKeeper supply.Keeper, distrKeeper distr.Keeper) {

	ctrl := controller.NewController(crdb, cdc, txDec, accKeeper, stkKeeper, suppKeeper, distrKeeper)
	blockchain.Init(ctrl, ctxCreator)
}

func InitGobGenesisState(cdc *codec.Codec, delegations staking.Delegations,
	validators staking.Validators, gs simapp.GenesisState) {
	prepareDB(crdb.GetGormClient())

	coinDenom := config.GetConfig().GetString("coins.denom")
	stakeDenom := config.GetConfig().GetString("stake.denom")
	// init genesis state such as accounts, auth, mint, distribution,...
	// currently only handle accounts, include genesis wallets and genesis validators
	authDataBz := gs[auth.ModuleName]
	var authData auth.GenesisState
	err := cdc.UnmarshalJSON(authDataBz, &authData)
	if err != nil {
		_ = email.SendMail(err.Error())
		panic(err)
	}
	listAccounts := make([]entity.WalletAddress, 0)
	for _, v := range authData.Accounts {
		originalVesting := ""
		delegatedFree := ""
		delegatedVesting := ""
		var startTime, endTime int64
		coinAIOZ := ""
		coinStake := ""
		for _, c := range v.GetCoins() {
			if c.Denom == coinDenom {
				coinAIOZ = c.Amount.String()
			}
			if c.Denom == stakeDenom {
				coinStake = c.Amount.String()
			}
		}
		decAIOZ, _ := decimal.NewFromString(coinAIOZ)
		decStake, _ := decimal.NewFromString(coinStake)
		switch accType := v.(type) {
		case *types.ContinuousVestingAccount:
			originalVesting = accType.OriginalVesting.String()
			delegatedFree = accType.DelegatedFree.String()
			delegatedVesting = accType.DelegatedVesting.String()
			startTime = accType.StartTime
			endTime = accType.EndTime
		case *types.DelayedVestingAccount:
			originalVesting = accType.OriginalVesting.String()
			delegatedFree = accType.DelegatedFree.String()
			delegatedVesting = accType.DelegatedVesting.String()
			endTime = accType.EndTime
			startTime = 0
		}
		listAccounts = append(listAccounts, entity.WalletAddress{
			Address:          v.GetAddress().String(),
			Coins:            v.GetCoins().String(),
			AccountNumber:    v.GetAccountNumber(),
			Sequence:         v.GetSequence(),
			CoinAIOZ:         decAIOZ,
			CoinStake:        decStake,
			OriginalVesting:  originalVesting,
			DelegatedFree:    delegatedFree,
			DelegatedVesting: delegatedVesting,
			StartTime:        startTime,
			EndTime:          endTime,
		})
	}
	for _, v := range validators {
		consPubkey, _ := cmtypes.Bech32ifyPubKey(cmtypes.Bech32PubKeyTypeConsPub, v.ConsPubKey)
		if err := crdb.PrepareTransaction(func(tx *gorm.DB) error {
			return tx.Exec(`UPSERT INTO 
			validators(address,tokens,power,jailed,status,val_cons_addr,val_cons_pubkey) 
			VALUES(?,?,?,?,?,?,?)`,
				v.OperatorAddress.String(), v.Tokens.String(), v.GetConsensusPower(), v.Jailed, v.Status.String(),
				v.GetConsAddr().String(), consPubkey).Error
		}); err != nil {
			_ = email.SendMail(err.Error())
			panic(err)
		}
	}
	for _, v := range delegations {
		if err := crdb.PrepareTransaction(func(tx *gorm.DB) error {
			return tx.Exec(`UPSERT INTO 
			delegators(delegator_address,validator_address,shares) 
			VALUES(?,?,?)`,
				v.DelegatorAddress.String(), v.ValidatorAddress.String(), v.Shares.String()).Error
		}); err != nil {
			_ = email.SendMail(err.Error())
			panic(err)
		}

		if err := crdb.PrepareTransaction(func(tx *gorm.DB) error {
			dec, _ := decimal.NewFromString(v.Shares.String())
			return tx.Exec(`UPSERT INTO 
			stakes(delegator_address,shares,shares_dec) 
			VALUES(?,?,?)`,
				v.DelegatorAddress.String(), v.Shares.String(), dec).Error
		}); err != nil {
			_ = email.SendMail(err.Error())
			panic(err)
		}
	}

	for _, v := range listAccounts {
		if err := crdb.GetGormClient().Model(&entity.WalletAddress{}).Create(v).Error; err != nil {
			panic(err)
		}
	}

	if err := crdb.FlushAllTxns(); err != nil {
		_ = email.SendMail(err.Error())
		panic(err)
	}
}
