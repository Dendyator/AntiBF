package db

import (
	"context"
	"fmt"
	"os"

	"github.com/Dendyator/AntiBF/internal/logger" //nolint
	"github.com/go-redis/redis/v8"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	db          *sqlx.DB
	redisClient *redis.Client
	ctx         = context.Background()
	appLogger   *logger.Logger
)

func InitDB(dsn string, log *logger.Logger) {
	appLogger = log
	appLogger.Infof("Using DSN: %s", dsn)
	appLogger.Infof("Connecting to database with DATABASE_URL: %s", os.Getenv("DATABASE_URL"))
	var err error
	db, err = sqlx.Open("pgx", dsn)
	if err != nil {
		appLogger.Fatalf("Failed to connect to database: %v", err)
	}

	if err = db.Ping(); err == nil {
		appLogger.Info("Successfully connected to the database!")
	} else {
		appLogger.Fatalf("Failed to connect to database: %v", err)
	}
}

func CloseDB() {
	if err := db.Close(); err != nil {
		appLogger.Fatalf("Failed to close database: %v", err)
	}
}

func InitRedis(address string, log *logger.Logger) {
	appLogger = log
	redisClient = redis.NewClient(&redis.Options{
		Addr: address,
	})

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		appLogger.Fatalf("Failed to connect to Redis: %v", err)
	}
}

func CloseRedis() {
	if err := redisClient.Close(); err != nil {
		appLogger.Fatalf("Failed to close Redis: %v", err)
	}
}

func CheckInRedis(listType, ip string) (bool, error) {
	if redisClient == nil {
		return false, fmt.Errorf("Redis client is not initialized")
	}

	val, err := redisClient.SIsMember(ctx, listType, ip).Result()
	if err != nil {
		appLogger.Warnf("Error checking Redis: %v", err)
		return false, err
	}
	return val, nil
}

func UpdateRedis(listType, subnet string, add bool) bool {
	if redisClient == nil {
		appLogger.Warn("Redis client is not initialized")
		return false
	}

	var err error
	if add {
		err = redisClient.SAdd(ctx, listType, subnet).Err()
	} else {
		err = redisClient.SRem(ctx, listType, subnet).Err()
	}
	if err != nil {
		appLogger.Warnf("Error updating Redis: %v", err)
		return false
	}
	return true
}
