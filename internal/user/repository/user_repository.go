package repository

import (
	"context"
	"math"
	"usermanagement-api/internal/user/domain"

	"gorm.io/gorm"
)

type (
	userRepository struct {
		db *gorm.DB
	}
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) RegisterUser(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetAllUserWithPagination(ctx context.Context, tx *gorm.DB, req PaginationRequest) (GetAllUserRepositoryResponse, error) {
	if tx == nil {
		tx = r.db
	}

	var users []UserModel
	var err error
	var count int64

	if req.PerPage == 0 {
		req.PerPage = 10
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if err := tx.WithContext(ctx).Model(&UserModel{}).Count(&count).Error; err != nil {
		return GetAllUserRepositoryResponse{}, err
	}

	if err := tx.WithContext(ctx).Scopes(Paginate(req.Page, req.PerPage)).Find(&users).Error; err != nil {
		return GetAllUserRepositoryResponse{}, err
	}

	totalPage := int64(math.Ceil(float64(count) / float64(req.PerPage)))

	return GetAllUserRepositoryResponse{
		Users: users,
		PaginationResponse: PaginationResponse{
			Page:    req.Page,
			PerPage: req.PerPage,
			Count:   count,
			MaxPage: totalPage,
		},
	}, err
}

func (r *userRepository) GetUserById(ctx context.Context, tx *gorm.DB, userId string) (User, error) {
	if tx == nil {
		tx = r.db
	}

	var user UserModel
	if err := tx.WithContext(ctx).Where("id = ?", userId).Take(&user).Error; err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (domain.User, error) {
	if tx == nil {
		tx = r.db
	}

	var user UserModel
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *userRepository) CheckEmail(ctx context.Context, tx *gorm.DB, email string) (domain.User, bool, error) {
	if tx == nil {
		tx = r.db
	}

	var user UserModel
	if err := tx.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		return domain.User{}, false, err
	}

	return user, true, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user domain.User) (domain.User, error) {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Updates(&user).Error; err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, tx *gorm.DB, userId string) error {
	if tx == nil {
		tx = r.db
	}

	if err := tx.WithContext(ctx).Delete(&UserModel{}, "id = ?", userId).Error; err != nil {
		return err
	}

	return nil
}
