package storage

import (
	"go-cars/storage/models"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func New(logger *zap.Logger, dsn string) *Storage {
	store := &Storage{
		Logger: logger,
	}

	store.init(dsn)

	return store
}

func (s *Storage) init(dsn string) {
	s.Logger.Debug("Connecting to the database...", zap.String("dsn: ", dsn))
	var db *gorm.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		s.Logger.Error("Failed to connect to the database", zap.String("Error: ", err.Error()))
		panic(err)
	}

	s.DB = db

	err = s.DB.AutoMigrate(&models.Cars{}, &models.Owners{})
	if err != nil {
		s.Logger.Error("Failed to migrate the database", zap.String("Error: ", err.Error()))
		return
	}
}
