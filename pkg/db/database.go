package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Yu-Qi/GoAuth/pkg/config"
	customlog "github.com/Yu-Qi/GoAuth/pkg/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

var database *gorm.DB
var initdatabaseOnce sync.Once

// Get database
func Get() *gorm.DB {
	initdatabaseOnce.Do(initialize)
	return database
}

// GetWith get database with context
func GetWith(ctx context.Context) *gorm.DB {
	initdatabaseOnce.Do(initialize)
	return database.WithContext(ctx)
}

// Init database
func Init() {
	initdatabaseOnce.Do(initialize)
}

// initialize will create a new database sesssion. If we are in an CI environment, a random table name will be used.
func initialize() {
	var err error
	username := config.GetString("MYSQL_USERNAME")
	password := config.GetString("MYSQL_PASSWORD")
	host := config.GetString("MYSQL_HOST")
	port := config.GetString("MYSQL_PORT")
	options := config.GetString("MYSQL_OPTIONS")
	databaseName := config.GetString("MYSQL_DATABASE")
	slowThreshold := time.Duration(config.GetInt("MYSQL_SLOW_THRESHOLD")) * time.Millisecond
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", username, password, host, port, databaseName, options)

	database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             slowThreshold,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		customlog.Fatalf("failed to connect database: %v", err)
	}
	if err = database.Use(tracing.NewPlugin()); err != nil {
		customlog.Fatalf("failed to use tracing plugin: %v", err)
	}

}
