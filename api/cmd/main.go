package main

import (
	"example/log"
	"go.uber.org/zap"
)

func main() {
	logger := log.Default()
	server := NewServer(logger)
	err := server.Start()

	logger.Error("server stopped", zap.Error(err))
}
