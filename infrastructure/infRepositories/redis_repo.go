package infRepositories

import (
	"context"
	"fmt"
	"github.com/Dendyator/AntiBF/infrastructure/database"
	"github.com/Dendyator/AntiBF/internal/repository"
	"github.com/Dendyator/AntiBF/pkg/logger"
)

type RedisRepo struct {
	db     *database.DB
	logger *logger.Logger
}

func NewRedisRepo(db *database.DB, logger *logger.Logger) repository.Repository {
	return &RedisRepo{
		db:     db,
		logger: logger,
	}
}

func (r *RedisRepo) CheckInList(listType, key string) (bool, error) {
	if r.db == nil {
		return false, fmt.Errorf("DB client is not initialized")
	}
	val, err := r.db.CheckInRedis(listType, key)
	if err != nil {
		r.logger.Warnf("Error checking Redis: %v", err)
		return false, err
	}
	return val, nil
}

func (r *RedisRepo) UpdateList(listType, key string, add bool) bool {
	var err error
	if add {
		err = r.db.Redis.SAdd(context.Background(), listType, key).Err()
	} else {
		err = r.db.Redis.SRem(context.Background(), listType, key).Err()
	}
	if err != nil {
		r.logger.Warnf("Error updating Redis: %v", err)
		return false
	}
	return true
}

func (r *RedisRepo) ClearBucket(key string) {
	r.db.Redis.Del(context.Background(), key)
}
