package main

import (
	"fmt"
	"go-cars/server"
	"go-cars/storage"
	"go-cars/utilities"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {
	godotenv.Load(".env")
	if !utilities.CheckEnv([]string{"DB_ADDR", "SRC_ADDR"}) {
		panic("Missing required ENVs")
	}
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println("Failed to create logger instance: ", err.Error())
	}

	logger.Info("Starting go-cars...")

	storage := storage.New(logger, os.Getenv("DB_ADDR"))

	logger.Debug("Creating server instance...")
	serv := server.New(storage, logger)

	serv.Run(":5000")
}
