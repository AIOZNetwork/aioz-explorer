package main

import (
	"aioz.io/go-aioz/app"
	"aioz.io/go-aioz/x_gob_explorer/config"
	"errors"
	"fmt"
	"github.com/tendermint/iavl"
	dbm "github.com/tendermint/tm-db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/tendermint/tendermint/proxy"
	tmsm "github.com/tendermint/tendermint/state"
)

var (
	conn *gorm.DB
)

var (
	home        string
	force       bool
	replayBlock int64
	host        string
	port        string
	user        string
	pwd         string
	dbname      string
	sslmode     string
	sslrootcert string
	sslcert     string
	sslkey      string
)

func main() {

	rootCmd := &cobra.Command{
		Use:   "replay",
		Short: "Replay aioz block",
		Run: func(cmd *cobra.Command, args []string) {
			initDB(host, port, user, pwd, dbname, sslmode, sslrootcert, sslcert, sslkey)
			runReplay(home)
		},
	}

	rootCmd.PersistentFlags().StringVar(&home, "home", "", "root dir")
	rootCmd.PersistentFlags().BoolVar(&force, "force", false, "Force delete all exceeded tree")
	rootCmd.PersistentFlags().Int64Var(&replayBlock, "replay", 0, "replay block number")
	host = config.GetConfig().GetString("database.host")
	port = config.GetConfig().GetString("database.port")
	user = config.GetConfig().GetString("database.user")
	pwd = config.GetConfig().GetString("database.passwd")
	dbname = config.GetConfig().GetString("database.dbname")
	sslmode = config.GetConfig().GetString("database.sslmode")
	sslrootcert = config.GetConfig().GetString("database.sslrootcert")
	sslcert = config.GetConfig().GetString("database.sslcert")
	sslkey = config.GetConfig().GetString("database.sslkey")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initDB(host, port, user, passwd, dbname, sslmode, sslrootcert, sslcert, sslkey string) {
	url := ""
	if sslmode == "disable" {
		url = fmt.Sprintf("user=%v host=%v port=%v dbname=%v",
			user, host, port, dbname)
	} else if sslmode == "require" {
		url = fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=%v sslrootcert=%v sslcert=%v sslkey=%v",
			user, passwd, host, port, dbname, sslmode, sslrootcert, sslcert, sslkey)
	} else {
		panic(errors.New("sslmode is undefined"))
	}
	c, err := gorm.Open(postgres.Open(url), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
	})
	if err != nil {
		panic(err)
	}

	conn = c
}

func runReplay(rootDir string) {
	var replay int64
	var err error
	if force {
		replay = replayBlock - 1
		if err := removeEntitiesToReplay(replayBlock); err != nil {
			panic(err)
		}
	} else {
		replay, err = getBlockReplay()
		if err != nil {
			panic(err)
		}
		replay -= 1
	}

	// Init node dependencies
	viper.Set("home", rootDir)

	//--------------------------------------
	dataDir := filepath.Join(rootDir, "data")
	appDB, tmDB, bcDB := loadDatabases(dataDir)

	// Application
	fmt.Println("Creating application")
	ctx := server.NewDefaultContext()
	myApp := app.NewAIOZApp(
		ctx.Logger, appDB, nil, true, 0,
		baseapp.SetPruning(store.PruneNothing),
	)

	cc := proxy.NewLocalClientCreator(myApp)
	proxyApp := proxy.NewAppConns(cc)
	err = proxyApp.Start()
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = proxyApp.Stop()
	}()

	// clear latest data
	if force {
		fmt.Println("Force clear block data")
		keys := myApp.GetKeys()

		for _, key := range keys {
			keyName := key.Name()

			db := dbm.NewPrefixDB(appDB, []byte("s/k:"+keyName+"/"))
			tree, err := iavl.NewMutableTree(db, 500000)
			if err != nil {
				panic(err)
			}
			_, err = tree.LoadVersion(replay)
			if err != nil {
				panic(err)
			}

			ver := replay + 1
			latestVersion := tree.GetLatestVersion()
			for ; ver <= latestVersion; ver++ {
				tree.ForceDeleteVersion(ver)
				if ver%1000 == 0 {
					fmt.Printf("Reset %v to version: %v\n", keyName, ver)
				}
			}

			tree.ForceDeleteCommit(replay)
			fmt.Printf("Reset %v to version: %v\n", keyName, replayBlock)
		}
	} else {
		err = myApp.LoadHeight(replay)
		if err != nil {
			panic(err)
		}
	}

	// Need to update atomically.
	batch := appDB.NewBatch()
	rootmulti.SetLatestVersion(batch, replay)
	batch.Write()

	// save state
	fmt.Println("Load and save state")
	state := loadState(replay, tmDB, bcDB)
	tmsm.SaveState(tmDB, state)

	//--------------------------------------------------
	// clear wal
	err = os.RemoveAll(filepath.Join(dataDir, "cs.wal"))
	if err != nil {
		panic(err)
	}

	// replay validator state
	varState := filepath.Join(dataDir, "priv_validator_state.json")
	read, err := ioutil.ReadFile(varState)
	if err != nil {
		panic(err)
	}

	r, _ := regexp.Compile(`"height": "(\d+)",`)
	readStr := r.ReplaceAllString(string(read), fmt.Sprintf(`"height": "%d",`, replay))
	readStr = strings.Replace(readStr, `"step": 3,`, `"step": 0,`, 1)

	err = ioutil.WriteFile(varState, []byte(readStr), 0644)
	if err != nil {
		panic(err)
	}

	// close
	appDB.Close()
	tmDB.Close()
	bcDB.Close()
}
