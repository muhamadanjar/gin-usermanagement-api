package repositories

import (
	"usermanagement-api/domain/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *entities.User) error
	FindByID(id uuid.UUID) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	FindByUsername(username string) (*entities.User, error)
	FindAll(page, pageSize int) ([]*entities.User, int64, error)
	Update(user *entities.User) error
	Delete(id uuid.UUID) error
	AssignRoles(userID uuid.UUID, roleIDs []uuid.UUID) error
	AppendRole(userID uuid.UUID, roleID uuid.UUID) error
	AddTokenHistory(history *entities.TokenHistory) error
	FindTokenHistory(userID uuid.UUID) ([]*entities.TokenHistory, error)
	RemoveTokenHistoryByToken(token string) error
}
