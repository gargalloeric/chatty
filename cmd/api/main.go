package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/gargalloeric/chatty/internal/chat"
	"github.com/gorilla/websocket"
)

type config struct {
	host string
	port string
	env  string
}

type application struct {
	logger   *slog.Logger
	config   config
	upgrader websocket.Upgrader
	hub      *chat.Hub
}

func main() {
	var conf config

	flag.StringVar(&conf.host, "host", "", "Server host address")
	flag.StringVar(&conf.port, "port", "3000", "Server port address")
	flag.StringVar(&conf.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	logger := slog.Default()
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	app := &application{
		logger:   logger,
		config:   conf,
		upgrader: upgrader,
		hub:      chat.NewHub(),
	}

	go app.hub.Run()

	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
