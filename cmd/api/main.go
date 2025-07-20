package main

import (
	"flag"
	"log/slog"
	"net"
	"os"
)

type config struct {
	host string
	port string
	env  string
}

type application struct {
	logger *slog.Logger
	config config
}

func main() {
	var conf config

	flag.StringVar(&conf.host, "host", "", "Server host address")
	flag.StringVar(&conf.port, "port", "3000", "Server port address")
	flag.StringVar(&conf.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	logger := slog.Default()

	app := &application{
		logger: logger,
		config: conf,
	}

	logger.Info("server listening", "addr", net.JoinHostPort(conf.host, conf.port))

	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
