package database

import (
	"context"
	"fmt"
	"github.com/Dendyator/AntiBF/pkg/logger"
	"github.com/go-redis/redis/v8"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	SQLDB     *sqlx.DB
	Redis     *redis.Client
	ctx       context.Context
	appLogger *logger.Logger
}

func NewDB(dsn, redisAddress string, log *logger.Logger) *DB {
	db := &DB{
		appLogger: log,
		ctx:       context.Background(),
	}

	sqlDB, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err = sqlDB.Ping(); err == nil {
		log.Info("Successfully connected to the database!")
	} else {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.SQLDB = sqlDB

	// Инициализация Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	db.Redis = redisClient

	return db
}

func (d *DB) Close() {
	if d.SQLDB != nil {
		if err := d.SQLDB.Close(); err != nil {
			d.appLogger.Errorf("Failed to close SQL database: %v", err)
		}
	}
	if d.Redis != nil {
		if err := d.Redis.Close(); err != nil {
			d.appLogger.Errorf("Failed to close Redis: %v", err)
		}
	}
}

// CheckInRedis проверяет наличие ключа в Redis.
func (d *DB) CheckInRedis(listType, ip string) (bool, error) {
	if d.Redis == nil {
		return false, fmt.Errorf("Redis client is not initialized")
	}
	val, err := d.Redis.SIsMember(d.ctx, listType, ip).Result()
	if err != nil {
		d.appLogger.Warnf("Error checking Redis: %v", err)
		return false, err
	}
	return val, nil
}

// UpdateRedis обновляет список в Redis.
func (d *DB) UpdateRedis(listType, subnet string, add bool) bool {
	if d.Redis == nil {
		d.appLogger.Warn("Redis client is not initialized")
		return false
	}
	var err error
	if add {
		err = d.Redis.SAdd(d.ctx, listType, subnet).Err()
	} else {
		err = d.Redis.SRem(d.ctx, listType, subnet).Err()
	}
	if err != nil {
		d.appLogger.Warnf("Error updating Redis: %v", err)
		return false
	}
	return true
}
