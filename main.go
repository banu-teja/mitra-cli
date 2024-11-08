package main

import (
	"database/sql"

	"github.com/rs/zerolog/log"

	"github.com/banu-teja/mitra-cli/cmd"
	db "github.com/banu-teja/mitra-cli/db/sqlc"
	"github.com/banu-teja/mitra-cli/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	connPool, err := sql.Open("sqlite3", "./history.db")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db:")
	}

	store := db.NewStore(connPool)

	cmd.Execute(config, store)
}
