package usecase

import (
	"github.com/Dendyator/AntiBF/internal/entity"
	"github.com/Dendyator/AntiBF/pkg/config"
	"github.com/Dendyator/AntiBF/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

// MockRepository - мок репозитория для тестирования.
type MockRepository struct {
}

func (m *MockRepository) CheckInList(listType, key string) (bool, error) {
	switch listType {
	case "whitelist":
		return key == "192.168.1.10/25", nil
	case "blacklist":
		return key == "192.168.1.20/25", nil
	default:
		return false, nil
	}
}

func (m *MockRepository) UpdateList(listType, key string, add bool) bool {
	return true // Для тестов всегда возвращаем успешное обновление
}

func (m *MockRepository) ClearBucket(key string) {
	// Ничего не делаем в моке
}

// MockUserRepository - мок userRepository (если он используется).
type MockUserRepository struct{}

func (m *MockUserRepository) GetUserByLogin(login string) (*entity.User, error) {
	return &entity.User{Login: login, Password: "testpass"}, nil
}

func TestCheckAuthorization(t *testing.T) {
	appLogger := logger.New("info")
	mockRepo := &MockRepository{}
	mockUserRepo := &MockUserRepository{}
	rateLimiter := NewRateLimiter(mockRepo, mockUserRepo, config.RateLimiterConfig{
		LoginLimit:    10,
		PasswordLimit: 100,
		IPLimit:       1000,
	}, appLogger)

	// Тест для whitelisted IP
	ok := rateLimiter.CheckAuthorization("testuser", "testpass", "192.168.1.10/25")
	assert.True(t, ok, "Expected true for whitelisted IP, got false")

	// Тест для blacklisted IP
	ok = rateLimiter.CheckAuthorization("testuser", "testpass", "192.168.1.20/25")
	assert.False(t, ok, "Expected false for blacklisted IP, got true")

	// Тест для валидных попыток авторизации
	for i := 0; i < 10; i++ {
		ok = rateLimiter.CheckAuthorization("testuser", "testpass", "192.168.1.30/25")
		if !ok {
			t.Errorf("Expected true for valid login attempt, got false at iteration %d", i)
		}
	}

	// Проверка превышения лимита
	ok = rateLimiter.CheckAuthorization("testuser", "testpass", "192.168.1.30/25")
	assert.False(t, ok, "Expected false after exceeding rate limit, got true")
}

func TestResetBucket(t *testing.T) {
	appLogger := logger.New("info")
	mockRepo := &MockRepository{}
	mockUserRepo := &MockUserRepository{}
	rateLimiter := NewRateLimiter(mockRepo, mockUserRepo, config.RateLimiterConfig{
		LoginLimit:    10,
		PasswordLimit: 100,
		IPLimit:       1000,
	}, appLogger)

	// Выполняем несколько попыток авторизации для создания bucket'ов
	for i := 0; i < 5; i++ {
		rateLimiter.CheckAuthorization("testuser", "testpass", "192.168.1.30/25")
	}

	// Сбрасываем bucket
	rateLimiter.ResetBucket("testuser", "192.168.1.30/25")

	// Проверяем, что bucket сброшен
	for _, key := range []string{"testuser", "192.168.1.30/25"} {
		if _, exists := rateLimiter.limiter.Load(key); exists {
			t.Errorf("Expected bucket to be reset for key %s, but entry still exists", key)
		}
	}
}

func TestAuthorizationLimits(t *testing.T) {
	appLogger := logger.New("debug")
	mockRepo := &MockRepository{}
	mockUserRepo := &MockUserRepository{}
	rateLimiter := NewRateLimiter(mockRepo, mockUserRepo, config.RateLimiterConfig{
		LoginLimit:    5,
		PasswordLimit: 5,
		IPLimit:       5,
	}, appLogger)

	// Сброс bucket'ов перед тестом
	rateLimiter.ResetBucket("testuser", "192.168.1.1/25")
	rateLimiter.ResetBucket("testpass", "192.168.1.1/25")

	// Проверяем лимиты
	for i := 0; i < 5; i++ {
		ok := rateLimiter.CheckAuthorization("testuser", "testpass", "192.168.1.1/25")
		assert.True(t, ok, "Authorization should be allowed")
	}

	// Превышение лимита
	ok := rateLimiter.CheckAuthorization("testuser", "testpass", "192.168.1.1/25")
	assert.False(t, ok, "Authorization should be blocked after exceeding limit")
}
