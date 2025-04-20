package main

import (
	"context"
	"os"
	"p2platform/api"
	db "p2platform/db/sqlc"
	"p2platform/util"
	"p2platform/worker"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "p2platform/docs"
)

//	@title			P2Platform API
//	@description	This is a simple API for a P2P platform.
//	@description
//	@description	Feel free to contact me if you have any questions
//	@description
//	@description				GitHub Repository:
//	@contact.name				Javakhir Yu
//	@contact.url				https://github.com/javakhiryu/p2platform
//	@contact.email				javakhiryulchibaev@gmail.com
//	@host						localhost:8080
//	@version					1.0
//	@BasePath					/
//	@schemes					http
//	@produce					json
//	@consumes					json

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to database")
	}
	store := db.NewStore(conn)
	worker := worker.NewAutoReleaseWorker(store, 1*time.Minute)
	worker.Start()
	runGinServer(config, store)

	defer worker.Stop()

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(store, config)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
		return
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start the server")
	}
}
