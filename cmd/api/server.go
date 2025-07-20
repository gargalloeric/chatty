package main

import (
	"log/slog"
	"net"
	"net/http"
	"time"
)

func (app *application) serve() error {
	server := &http.Server{
		Addr:         net.JoinHostPort(app.config.host, app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	return server.ListenAndServe()
}
