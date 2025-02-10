package repositories

import (
	"github.com/Dendyator/AntiBF/infrastructure/database"
	"github.com/Dendyator/AntiBF/internal/entity"
	"github.com/Dendyator/AntiBF/internal/repository"
	"github.com/Dendyator/AntiBF/pkg/logger"
)

// UserRepository представляет реализацию интерфейса UserRepositoryInterface.
type UserRepository struct {
	db     *database.DB // Ссылка на экземпляр базы данных
	logger *logger.Logger
}

// NewUserRepository создаёт новый экземпляр UserRepository.
func NewUserRepository(db *database.DB, logger *logger.Logger) repository.UserRepositoryInterface {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// GetUserByLogin получает пользователя из базы данных по логину.
func (r *UserRepository) GetUserByLogin(login string) (*entity.User, error) {
	var user entity.User
	err := r.db.SQLDB.Get(&user, "SELECT id, login, password FROM users WHERE login = $1", login)
	if err != nil {
		r.logger.Warnf("Failed to get user by login %s: %v", login, err)
		return nil, err
	}
	return &user, nil
}
