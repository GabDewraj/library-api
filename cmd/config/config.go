package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ServerPort       string
	DB               DBConfig
	RedisConfig      RedisConfig
	MiddlewareConfig MiddlewareConfig
}

// Mysql DB config
type DBConfig struct {
	Driver                 string
	Host                   string
	Port                   string
	Database               string
	Password               string
	Username               string
	MigrationDirectoryPath string
	ForceTLS               bool
}

// Redis DB config
type RedisConfig struct {
	Host string
	Port int
}
type MiddlewareConfig struct {
	MaxRequests int
	RateWindow  int
}

func NewConfig() (*Config, error) {
	// Retrieve params for rate limiting
	maxrequests, err := strconv.Atoi(os.Getenv("RATE_LIMITER_MAX_REQUESTS"))
	if err != nil {
		return nil, err
	}
	ratewindow, err := strconv.Atoi(os.Getenv("RATE_LIMITER_WINDOW"))
	if err != nil {
		return nil, err
	}
	// Create App Config Object from env
	redisport, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return nil, err
	}
	return &Config{
		ServerPort: fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		DB: DBConfig{
			Driver:                 os.Getenv("MYSQL_DRIVER"),
			Host:                   os.Getenv("MYSQL_HOST"),
			Port:                   os.Getenv("MYSQL_PORT"),
			Database:               os.Getenv("MYSQL_DATABASE"),
			Password:               os.Getenv("MYSQL_PASSWORD"),
			Username:               os.Getenv("MYSQL_USERNAME"),
			MigrationDirectoryPath: os.Getenv("SERVER_MIGRATION_DIRECTORY"),
			ForceTLS:               false,
		},
		RedisConfig: RedisConfig{
			Host: os.Getenv("REDIS_HOST"),
			Port: redisport,
		},
		MiddlewareConfig: MiddlewareConfig{
			MaxRequests: maxrequests,
			RateWindow:  ratewindow,
		},
	}, nil
}

// Database Connection Configuration.
func NewDBConnection(config *Config) (*sqlx.DB, error) {
	dbConfig := config.DB
	logger := logrus.StandardLogger()
	logger.Infoln("Connecting to MySQL DB")
	logger.WithFields(logrus.Fields{
		"username": dbConfig.Username,
		"host":     dbConfig.Host,
		"port":     dbConfig.Port,
		"database": dbConfig.Database,
	}).Debug("connecting to db")
	dbAddress := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&tls=%t",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
		dbConfig.ForceTLS)
	// If running in docker compose allow for db to initialise
	var db *sqlx.DB
	for retries := 0; retries < 100; retries++ {
		successConn, err := sqlx.Connect("mysql", dbAddress)
		if err == nil {
			// Connection successful, break out of the loop
			db = successConn
			break
		}

		// Print or log the error (optional)
		fmt.Printf("Failed to connect to the database. Retrying... (Attempt %d)\n", retries+1)

		// Wait for the specified interval before retrying
		time.Sleep(4 * time.Second)
	}

	logger.Infoln("Successfully Connected to Database Host")
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}

func PerformMigrations(db *sqlx.DB, logger *logrus.Logger, args *Config) error {
	logger.Infoln("Performing migrations")
	n, err := migrate.Exec(db.DB, "mysql", &migrate.FileMigrationSource{
		Dir: args.DB.MigrationDirectoryPath,
	}, migrate.Up)
	logger.Infof("Performed %d migrations", n)
	return err
}

// Create a Redis client for Cache service
func NewRedisClient(config *Config) (*redis.Client, error) {
	redisConfig := config.RedisConfig
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
