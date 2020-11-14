package main

import (
	"errors"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"users-api/api/http"
	"users-api/models"
	"users-api/storage"
)

const (
	storageDriver = "STORAGE_DRIVER"
	dbUsernameEnv = "DATABASE_USERNAME"
	dbPasswordEnv = "DATABASE_PASSWORD"
	dbNameEnv     = "DATABASE_NAME"
	dbHostEnv     = "DATABASE_HOST"
	dbPortEnv     = "DATABASE_PORT"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Error loading .env file")
	}

	userStorage, err := getStorageDriver()
	if err != nil {
		logrus.WithError(err).Fatal("Storage driver initializing")
	}

	server := http.Server{Storage: userStorage, Port: "8080"}

	logrus.Info("Starting rest server")
	logrus.WithError(server.Start()).Fatal("Rest server can't start")
}

func getStorageDriver() (driver models.UserStorage, err error) {
	switch os.Getenv(storageDriver) {
	case "inmemory":
		driver = storage.NewInMemoryStorage()
	case "postgresql":
		port, err := strconv.Atoi(os.Getenv(dbPortEnv))
		if err != nil {
			logrus.WithError(err).Fatal("Converting port to int")
		}

		driver = storage.NewPostgreSqlStorage(storage.Config{
			Username:     os.Getenv(dbUsernameEnv),
			Password:     os.Getenv(dbPasswordEnv),
			DatabaseName: os.Getenv(dbNameEnv),
			Host:         os.Getenv(dbHostEnv),
			Port:         port,
		})
	default:
		err = errors.New("unknown database driver, check your .env file")
	}
	return
}
