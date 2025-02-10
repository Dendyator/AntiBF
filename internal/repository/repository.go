package repository

import "github.com/Dendyator/AntiBF/internal/entity"

type Repository interface {
	CheckInList(listType, key string) (bool, error)
	UpdateList(listType, key string, add bool) bool
	ClearBucket(key string)
}

type UserRepositoryInterface interface {
	GetUserByLogin(login string) (*entity.User, error)
}
