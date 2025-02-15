package main

import (
	"embed"

	"github.com/rs/zerolog/log"

	"github.com/banu-teja/mitra-cli/cmd"
	db "github.com/banu-teja/mitra-cli/db/sqlc"
	"github.com/banu-teja/mitra-cli/utils"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed db/migration/*
var migrationFiles embed.FS

func main() {
	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	utils.MigrationFiles = migrationFiles

	sqlDB, err := utils.InitializeDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database:")
	}
	defer sqlDB.Close() // Ensure the database is closed when main exits

	// connPool, err := sql.Open("sqlite3", "~/.ai/history.db")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("cannot connect to db:")
	// }

	store := db.NewStore(sqlDB)

	cmd.Execute(config, store)
}
