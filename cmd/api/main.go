package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"sync"

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
	room     *chat.Room
	wg       sync.WaitGroup
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
		room:     chat.NewRoom(context.Background(), logger, "Test Room"),
	}

	app.background(func() {
		app.room.Run()
	})

	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
