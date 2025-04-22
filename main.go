package main

import (
	"context"
	"embed"
	"net/http"
	"os"
	"p2platform/api"
	db "p2platform/db/sqlc"
	"p2platform/notification/kafka"
	"p2platform/notification/telegram"
	"p2platform/util"
	"p2platform/worker"
	"strings"
	"time"

	_ "p2platform/docs"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

//go:embed static/*
var staticFS embed.FS

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

	tg := telegram.New(config.TelegramBotToken)

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				data, _ := staticFS.ReadFile("static/index.html")
				w.Write(data)
				return
			}
			http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
		})

		log.Info().Msg("🌐 Serving on http://localhost:8081")
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			log.Error().Err(err).Msg("HTTPS static server failed")
		}
	}()

	go func() {
		log.Info().Msg("Starting Telegram Kafka consumer...")
		err := kafka.StartConsumer(strings.Split(config.KafkaBrokers, ","), "notifications", tg)
		if err != nil {
			log.Error().Err(err).Msg("Kafka consumer error")
		}
	}()

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
