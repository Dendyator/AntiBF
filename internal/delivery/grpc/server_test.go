package grpc

import (
	"context"
	"github.com/Dendyator/AntiBF/internal/delivery/grpc/proto/pb"
	"github.com/Dendyator/AntiBF/internal/repositories"
	"github.com/Dendyator/AntiBF/pkg/config"
	"github.com/Dendyator/AntiBF/pkg/logger"
	"testing"

	"github.com/Dendyator/AntiBF/internal/usecase" //nolint
	"github.com/stretchr/testify/assert"
)

func TestCheckAuthorization(t *testing.T) {
	mockLogger := logger.New("info")

	// Создаем моковые репозитории
	mockRepo := &repositories.MockRepository{}
	mockUserRepo := &repositories.MockUserRepository{}

	// Создаем RateLimiter
	rateLimiter := usecase.NewRateLimiter(mockRepo, mockUserRepo, config.RateLimiterConfig{
		LoginLimit:    10,
		PasswordLimit: 100,
		IPLimit:       1000,
	}, mockLogger)

	// Создаем gRPC-сервер
	srv := NewServer(rateLimiter, mockLogger)

	// Подготавливаем запрос
	req := &pb.AuthRequest{
		Login:    "testuser",
		Password: "testpass",
		Ip:       "192.168.1.30/25",
	}
	resp, err := srv.CheckAuthorization(context.Background(), req)

	// Проверяем ошибку и ответ
	assert.NoError(t, err)
	assert.True(t, resp.Ok)
}

func TestResetBucket(t *testing.T) {
	mockLogger := logger.New("info")

	// Создаем моковые репозитории
	mockRepo := &repositories.MockRepository{}
	mockUserRepo := &repositories.MockUserRepository{}

	// Создаем RateLimiter
	rateLimiter := usecase.NewRateLimiter(mockRepo, mockUserRepo, config.RateLimiterConfig{
		LoginLimit:    10,
		PasswordLimit: 100,
		IPLimit:       1000,
	}, mockLogger)

	// Создаем gRPC-сервер
	srv := NewServer(rateLimiter, mockLogger)

	// Подготавливаем запрос
	req := &pb.ResetRequest{
		Login: "testuser",
		Ip:    "192.168.1.30/25",
	}
	resp, err := srv.ResetBucket(context.Background(), req)

	// Проверяем ошибку и ответ
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestAddToBlacklist(t *testing.T) {
	mockLogger := logger.New("info")

	// Создаем моковые репозитории
	mockRepo := &repositories.MockRepository{}
	mockUserRepo := &repositories.MockUserRepository{}

	// Создаем RateLimiter
	rateLimiter := usecase.NewRateLimiter(mockRepo, mockUserRepo, config.RateLimiterConfig{
		LoginLimit:    10,
		PasswordLimit: 100,
		IPLimit:       1000,
	}, mockLogger)

	// Создаем gRPC-сервер
	srv := NewServer(rateLimiter, mockLogger)

	// Подготавливаем запрос
	req := &pb.ListRequest{
		Subnet: "192.168.1.0/24",
	}
	resp, err := srv.AddToBlacklist(context.Background(), req)

	// Проверяем ошибку и ответ
	assert.NoError(t, err)
	assert.True(t, resp.Success)
}
