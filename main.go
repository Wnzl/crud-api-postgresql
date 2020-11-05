package main

import (
	"github.com/sirupsen/logrus"
	"users-api/api/http"
	"users-api/storage"
)

func main() {
	inMemoryStorage := storage.NewInMemoryStorage()

	server := http.Server{Storage: inMemoryStorage, Port: "8080"}

	logrus.Info("Starting rest server")
	logrus.WithError(server.Start()).Fatal("Rest server can't start")
}
