package x_gob_explorer

import (
	"aioz.io/go-aioz/x_gob_explorer/config"
	database "aioz.io/go-aioz/x_gob_explorer/domain/db"
	"aioz.io/go-aioz/x_gob_explorer/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	flagResetDB      = "resetDB"
	flagResetNode    = "resetData"

	flagWithTendermint = "with-tendermint"
	flagAddress        = "address"
	flagTraceStore     = "trace-store"
	flagPruning        = "pruning"
	flagCPUProfile     = "cpu-profile"
	FlagMinGasPrices   = "minimum-gas-prices"
)

func GobCmd(start func() error) *cobra.Command {

	host := config.GetConfig().GetString("database.host")
	port := config.GetConfig().GetString("database.port")
	user := config.GetConfig().GetString("database.user")
	pwd := config.GetConfig().GetString("database.passwd")
	dbname := config.GetConfig().GetString("database.dbname")
	sslmode := config.GetConfig().GetString("database.sslmode")
	sslrootcert := config.GetConfig().GetString("database.sslrootcert")
	sslkey := config.GetConfig().GetString("database.sslkey")
	sslcert := config.GetConfig().GetString("database.sslcert")
	schema := config.GetConfig().GetString("database.schema")

	cmd := &cobra.Command{
		Use:   "gob",
		Short: "Run the full node with indexing",
		RunE: func(cmd *cobra.Command, args []string) error {

			resetDB := viper.GetBool(flagResetDB)
			resetNode := viper.GetBool(flagResetNode)

			if resetDB {
				db := database.NewCockroachDB(user, pwd, host, port,
					dbname, sslmode, sslrootcert, sslkey, sslcert)
				db.Reset(schema)
				return nil
			}
			if resetNode {
				homeDir := viper.GetString(cli.HomeFlag)
				if err := utils.ResetNodeDB(homeDir); err != nil {
					return err
				}
				return nil
			}
			return start()
		},
	}

	// core flags for the ABCI application
	cmd.Flags().Bool(flagResetDB, false, "Reset cockroach database")
	cmd.Flags().Bool(flagResetNode, false, "Delete node directory")

	cmd.Flags().Bool(flagWithTendermint, true, "Run abci app embedded in-process with tendermint")
	cmd.Flags().String(flagAddress, "tcp://0.0.0.0:26658", "Listen address")
	cmd.Flags().String(flagTraceStore, "", "Enable KVStore tracing to an output file")
	cmd.Flags().String(flagPruning, "syncable", "Pruning strategy: syncable, nothing, everything")
	cmd.Flags().String(flagCPUProfile, "", "Enable CPU profiling and write to the provided file")
	cmd.Flags().String(FlagMinGasPrices, "", "Minimum gas prices to accept for transactions; All fees in a tx must meet this minimum (e.g. 0.01photino,0.0001stake)")

	// add support for all Tendermint-specific command line options
	tcmd.AddNodeFlags(cmd)
	return cmd
}
