package main

import (
	"log"

	"github.com/obliviousfrog/flighttracker/internal/server"
	"github.com/obliviousfrog/flighttracker/internal/tracker"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Failed to setup service logging.")
	}

	server := server.New(server.Config{
		Host:    "127.0.0.1",
		Port:    8080,
		Tracker: tracker.New(),
		Log:     zapLogger,
	})

	zapLogger.Info("Starting Flight Tracker Service.",
		zap.String("host", "127.0.0.1"),
		zap.Int("port", 8080),
	)
	if err := server.Start(); err != nil {
		zapLogger.Fatal("Failed to start service", zap.Error(err))
	}
}
