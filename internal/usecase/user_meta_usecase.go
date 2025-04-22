// internal/usecase/user_meta_usecase.go
package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/dto"
	"usermanagement-api/pkg/cache"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserMetaUseCase interface {
	CreateOrUpdate(userID uuid.UUID, key, value string) error
	GetByKey(userID uuid.UUID, key string) (*dto.UserMetaResponse, error)
	GetAllByUserID(userID uuid.UUID) (map[string]string, error)
	Delete(userID uuid.UUID, key string) error
}

type userMetaUseCase struct {
	userMetaRepo repositories.UserMetaRepository
	cache        cache.Cache
}

func NewUserMetaUseCase(userMetaRepo repositories.UserMetaRepository, cache cache.Cache) UserMetaUseCase {
	return &userMetaUseCase{
		userMetaRepo: userMetaRepo,
		cache:        cache,
	}
}

func (uc *userMetaUseCase) CreateOrUpdate(userID uuid.UUID, key, value string) error {
	// Try to find existing meta
	existingMeta, err := uc.userMetaRepo.FindByUserIDAndKey(userID, key)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if existingMeta != nil {
		// Update existing
		existingMeta.Value = value
		if err := uc.userMetaRepo.Update(existingMeta); err != nil {
			return err
		}
	} else {
		// Create new
		userMeta := &entities.UserMeta{
			Key:    key,
			Value:  value,
			UserID: userID,
		}
		if err := uc.userMetaRepo.Create(userMeta); err != nil {
			return err
		}
	}

	// Clear cache
	cacheKey := uc.getUserMetaCacheKey(userID)
	ctx := context.Background()
	_ = uc.cache.Delete(ctx, cacheKey)

	return nil
}

func (uc *userMetaUseCase) GetByKey(userID uuid.UUID, key string) (*dto.UserMetaResponse, error) {
	// Check cache first
	ctx := context.Background()
	cacheKey := uc.getUserMetaCacheKey(userID)

	cachedData, err := uc.cache.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		// Parse cached data
		var metaMap map[string]string
		if err := json.Unmarshal([]byte(cachedData), &metaMap); err == nil {
			if value, ok := metaMap[key]; ok {
				return &dto.UserMetaResponse{
					Key:   key,
					Value: value,
				}, nil
			}
		}
	}

	// If not in cache, get from database
	userMeta, err := uc.userMetaRepo.FindByUserIDAndKey(userID, key)
	if err != nil {
		return nil, err
	}

	// Update cache with all user meta
	go uc.updateUserMetaCache(userID)

	return &dto.UserMetaResponse{
		ID:     userMeta.ID,
		Key:    userMeta.Key,
		Value:  userMeta.Value,
		UserID: userMeta.UserID,
	}, nil
}

func (uc *userMetaUseCase) GetAllByUserID(userID uuid.UUID) (map[string]string, error) {
	// Check cache first
	ctx := context.Background()
	cacheKey := uc.getUserMetaCacheKey(userID)

	cachedData, err := uc.cache.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		var metaMap map[string]string
		if err := json.Unmarshal([]byte(cachedData), &metaMap); err == nil {
			return metaMap, nil
		}
	}

	// If not in cache, get from database
	metaMap, err := uc.userMetaRepo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Update cache
	_ = uc.cache.Set(ctx, cacheKey, metaMap, 30*time.Minute)

	return metaMap, nil
}

func (uc *userMetaUseCase) Delete(userID uuid.UUID, key string) error {
	userMeta, err := uc.userMetaRepo.FindByUserIDAndKey(userID, key)
	if err != nil {
		return err
	}

	if err := uc.userMetaRepo.Delete(userMeta.ID); err != nil {
		return err
	}

	// Clear cache
	cacheKey := uc.getUserMetaCacheKey(userID)
	ctx := context.Background()
	_ = uc.cache.Delete(ctx, cacheKey)

	return nil
}

func (uc *userMetaUseCase) getUserMetaCacheKey(userID uuid.UUID) string {
	return fmt.Sprintf("user_meta:%s", userID.String())
}

func (uc *userMetaUseCase) updateUserMetaCache(userID uuid.UUID) {
	ctx := context.Background()
	cacheKey := uc.getUserMetaCacheKey(userID)

	metaMap, err := uc.userMetaRepo.GetAllByUserID(userID)
	if err == nil {
		_ = uc.cache.Set(ctx, cacheKey, metaMap, 30*time.Minute)
	}
}
