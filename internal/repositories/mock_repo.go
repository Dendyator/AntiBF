package repositories

import (
	"errors"
	"github.com/Dendyator/AntiBF/internal/entity"
	"github.com/Dendyator/AntiBF/pkg/logger"
	"strconv"
)

type MockRepository struct{}

type MockUserRepository struct {
	logger *logger.Logger
}

func (m *MockRepository) CheckInList(listType, key string) (bool, error) {
	switch listType {
	case "whitelist":
		return key == "192.168.1.1/25", nil
	case "blacklist":
		return key == "192.168.1.2/25", nil
	default:
		return false, errors.New("invalid list type")
	}
}

func (m *MockRepository) UpdateList(listType, key string, add bool) bool {
	return true
}

func (m *MockRepository) ClearBucket(key string) {
	// Ничего не делаем
}

func NewMockUserRepository(log *logger.Logger) *MockUserRepository {
	return &MockUserRepository{
		logger: log,
	}
}

// GetUserByLogin возвращает пользователя по логину (моковая реализация).
func (m *MockUserRepository) GetUserByLogin(login string) (*entity.User, error) {
	m.logger.Debugf("MockUserRepository: Fetching user by login: %s", login)

	// Возвращаем фиктивного пользователя для тестирования
	if login == "testuser" {
		return &entity.User{
			ID:       strconv.Itoa(1),
			Login:    "testuser",
			Password: "testpass",
		}, nil
	}

	// Если пользователь не найден, возвращаем ошибку
	return nil, nil
}
