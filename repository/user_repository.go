package repository

import (
	"context"
	"math"
	"usermanagement-api/domain/models"

	"gorm.io/gorm"
)

type (
	UserRepository interface {
		RegisterUser(ctx context.Context, tx *gorm.DB, user models.User) (models.User, error)
		GetUserById(ctx context.Context, tx *gorm.DB, userId string) (models.User, error)
		GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (models.User, error)
		CheckEmail(ctx context.Context, tx *gorm.DB, email string) (models.User, bool, error)
		UpdateUser(ctx context.Context, tx *gorm.DB, user models.User) (models.User, error)
		DeleteUser(ctx context.Context, tx *gorm.DB, userId string) error
	}

	userRepository struct {
		db *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) RegisterUser(ctx context.Context, tx *gorm.DB, user models.User) (models.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req models.PaginationRequest) (models.GetAllUserRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var users []models.User
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if err := tx.WithContext(ctx).Model(&models.User{}).Count(&count).Error; err != nil {
		return models.GetAllUserRepositoryResponse{}, err
	}

	if err := tx.WithContext(ctx).Scopes(Paginate(req.Page, req.PerPage)).Find(&users).Error; err != nil {
		return models.GetAllUserRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return models.GetAllUserRepositoryResponse{
		Users: users,
		PaginationResponse: models.PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, err
}

func (r *userRepository) GetUserById(ctx context.Context, tx *gorm.DB, userId string) (models.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user models.User
	if err := tx.WithContext(ctx).Where("id = ?", userId).Take(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (models.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user models.User
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepository) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (models.User, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var user models.User
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return models.User{}, false, err
	}

	return user, true, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user models.User) (models.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Updates(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, tx *gorm.DB, userId string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&models.User{}, "id = ?", userId).Error; err != nil {
		return err
	}

	return nil
}
