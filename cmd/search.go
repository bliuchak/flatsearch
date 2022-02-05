package cmd

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/spf13/viper"

	"github.com/bliuchak/flatsearch/internal/estate"
	srealityplatform "github.com/bliuchak/flatsearch/internal/platform/sreality"
	"github.com/bliuchak/flatsearch/internal/platform/storage"
	"github.com/bliuchak/flatsearch/internal/sreality"
	"github.com/bliuchak/flatsearch/internal/sreality/decoder"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search an estate.",
	Run:   run,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().String("database_dsn", "", "Database DSN")
}

func run(cmd *cobra.Command, args []string) {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger := zerolog.New(os.Stderr).With().Str("cmd", "search").Timestamp().Logger()

	ctx := context.Background()

	logger.Info().Msg("start")

	sqlConn, err := sql.Open("postgres", viper.GetString("database_dsn"))
	if err != nil {
		logger.Fatal().Err(err).Msg("can't connect to postgres")
	}

	pgDB := sqlx.NewDb(sqlConn, "postgres")
	if err := pgDB.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("ping to postgres")
	}

	defer pgDB.Close()

	api := sreality.NewSrealityAPI(srealityplatform.NewClient(), &decoder.JSON{})
	resp, err := api.GetAll()
	if err != nil {
		panic(err)
	}

	logger.Info().Int("estates_total", len(resp.Estates)).Msg("get estates from api")

	handler := estate.NewEstate(api, storage.NewPostgres(pgDB), logger)
	err = handler.CompareAndSave(ctx, resp.Estates)
	if err != nil {
		panic(err)
	}

	logger.Info().Msg("finish")
}
