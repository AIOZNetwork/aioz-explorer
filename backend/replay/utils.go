package main

import (
	"aioz.io/go-aioz/types"
	"aioz.io/go-aioz/x_gob_explorer/domain/entity"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmsm "github.com/tendermint/tendermint/state"
	tmstore "github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/version"
	dbm "github.com/tendermint/tm-db"
)

func getBlockReplay() (int64, error) {
	latestBlock := new(entity.Block)
	if err := conn.Last(&latestBlock).Error; err != nil {
		panic(err)
	}
	return latestBlock.Height, removeEntitiesToReplay(latestBlock.Height)
}

func removeEntitiesToReplay(height int64) error {
	// delete block
	if err := conn.Where("height >= ?", height).Delete(&entity.Block{}).Error; err != nil {
		return err
	}

	// delete transaction
	if err := conn.Where("block_height >= ?", height).Delete(&entity.Transaction{}).Error; err != nil {
		return err
	}

	// delete messages and all kind of messages
	if err := conn.Where("block_height >= ?", height).Delete(&entity.Txs{}).Error; err != nil {
		return err
	}
	if err := conn.Where("block_height >= ?", height).Delete(&entity.MessageSend{}).Error; err != nil {
		return err
	}
	
	return nil
}

func initConfig() {
	// config prefix
	sdkConfig := sdk.GetConfig()
	sdkConfig.SetBech32PrefixForAccount(types.Bech32PrefixAccAddr, types.Bech32PrefixAccPub)
	sdkConfig.SetBech32PrefixForValidator(types.Bech32PrefixValAddr, types.Bech32PrefixValPub)
	sdkConfig.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr, types.Bech32PrefixConsPub)
	sdkConfig.SetCoinType(types.CoinType)
	sdkConfig.SetFullFundraiserPath(types.FullFundraiserPath)
	sdkConfig.Seal()
}

func loadDatabases(dataDir string) (dbm.DB, dbm.DB, dbm.DB) {
	// App DB
	fmt.Println("Opening app database")
	appDB, err := sdk.NewLevelDB("application", dataDir)
	if err != nil {
		panic(err)
	}

	// TM DB
	fmt.Println("Opening tendermint state database")
	tmDB, err := sdk.NewLevelDB("state", dataDir)
	if err != nil {
		panic(err)
	}

	// Blockchain DB
	fmt.Println("Opening blockstore database")
	bcDB, err := sdk.NewLevelDB("blockstore", dataDir)
	if err != nil {
		panic(err)
	}

	return appDB, tmDB, bcDB
}

func loadState(replayBlock int64, tmDB, bcDB dbm.DB) tmsm.State {
	// Create block store
	fmt.Printf("Loading block %v state\n", replayBlock)
	tmstore.BlockStoreStateJSON{Height: replayBlock}.Save(bcDB)
	blockStore := tmstore.NewBlockStore(bcDB)

	rollbackBlock1 := blockStore.LoadBlock(replayBlock - 1)
	rollbackBlock := blockStore.LoadBlock(replayBlock)

	consensusParams, err := tmsm.LoadConsensusParams(tmDB, replayBlock)
	if err != nil {
		panic(err)
	}

	validators, err := tmsm.LoadValidators(tmDB, replayBlock)
	if err != nil {
		panic(err)
	}

	nextValidators, err := tmsm.LoadValidators(tmDB, replayBlock+1)
	if err != nil {
		panic(err)
	}

	lastValidators, err := tmsm.LoadValidators(tmDB, replayBlock-1)
	if err != nil {
		panic(err)
	}

	state := tmsm.State{
		Version: tmsm.Version{
			Consensus: rollbackBlock.Version,
			Software:  version.TMCoreSemVer,
		},
		ChainID: rollbackBlock.ChainID,

		LastBlockHeight: rollbackBlock1.Header.Height,
		LastBlockTime:   rollbackBlock1.Header.Time,
		LastBlockID:     rollbackBlock.LastBlockID,
		LastResultsHash: rollbackBlock.LastResultsHash,

		NextValidators:              nextValidators,
		Validators:                  validators,
		LastValidators:              lastValidators,
		LastHeightValidatorsChanged: rollbackBlock1.Header.Height + 1 + 1,

		ConsensusParams:                  consensusParams,
		LastHeightConsensusParamsChanged: 1,

		AppHash: rollbackBlock.AppHash,
	}
	return state
}

func FindInt64(slice []int64, val int64) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
