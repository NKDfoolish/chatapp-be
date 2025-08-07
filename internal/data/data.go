package data

import (
	"kchat-be/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewGreeterRepo,
	NewUserRepo,
)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(logger)

	// Configure GORM logger
	gormLogger := gormlogger.Default.LogMode(gormlogger.Info)

	// Connect to database
	db, err := gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		helper.Errorf("failed to connect database: %v", err)
		return nil, nil, err
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		helper.Errorf("failed to get database instance: %v", err)
		return nil, nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		helper.Errorf("failed to ping database: %v", err)
		return nil, nil, err
	}

	helper.Info("database connected successfully")

	cleanup := func() {
		helper.Info("closing the data resources")
		sqlDB.Close()
	}

	return &Data{db: db}, cleanup, nil
}
