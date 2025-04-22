// internal/usecase/setting_usecase.go
package usecase

import (
	"context"
	"encoding/json"
	"time"
	"usermanagement-api/domain/entities"
	"usermanagement-api/domain/repositories"
	"usermanagement-api/internal/constants"
	"usermanagement-api/internal/dto"
	"usermanagement-api/pkg/cache"
)

type SettingUseCase interface {
	CreateOrUpdate(key, value string) error
	GetByKey(key string) (*dto.SettingResponse, error)
	GetAll() (map[string]string, error)
	Delete(key string) error
}

type settingUseCase struct {
	settingRepo repositories.SettingRepository
	cache       cache.Cache
}

func NewSettingUseCase(settingRepo repositories.SettingRepository, cache cache.Cache) SettingUseCase {
	return &settingUseCase{
		settingRepo: settingRepo,
		cache:       cache,
	}
}

func (uc *settingUseCase) CreateOrUpdate(key, value string) error {
	setting := &entities.Setting{
		Key:   key,
		Value: value,
	}

	if err := uc.settingRepo.Upsert(setting); err != nil {
		return err
	}

	// Clear cache
	ctx := context.Background()
	_ = uc.cache.Delete(ctx, constants.SettingsCacheKey)

	return nil
}

func (uc *settingUseCase) GetByKey(key string) (*dto.SettingResponse, error) {
	// Check cache first
	ctx := context.Background()
	cachedData, err := uc.cache.Get(ctx, constants.SettingsCacheKey)
	if err == nil && cachedData != "" {
		var settingsMap map[string]string
		if err := json.Unmarshal([]byte(cachedData), &settingsMap); err == nil {
			if value, ok := settingsMap[key]; ok {
				return &dto.SettingResponse{
					Key:   key,
					Value: value,
				}, nil
			}
		}
	}

	// If not in cache, get from database
	setting, err := uc.settingRepo.FindByKey(key)
	if err != nil {
		return nil, err
	}

	// Update cache with all settings
	go uc.updateSettingsCache()

	return &dto.SettingResponse{
		Key:   setting.Key,
		Value: setting.Value,
	}, nil
}

func (uc *settingUseCase) GetAll() (map[string]string, error) {
	// Check cache first
	ctx := context.Background()
	cachedData, err := uc.cache.Get(ctx, constants.SettingsCacheKey)
	if err == nil && cachedData != "" {
		var settingsMap map[string]string
		if err := json.Unmarshal([]byte(cachedData), &settingsMap); err == nil {
			return settingsMap, nil
		}
	}

	// If not in cache, get from database
	settings, err := uc.settingRepo.FindAll()
	if err != nil {
		return nil, err
	}

	settingsMap := make(map[string]string)
	for _, setting := range settings {
		settingsMap[setting.Key] = setting.Value
	}

	// Update cache
	_ = uc.cache.Set(ctx, constants.SettingsCacheKey, settingsMap, 1*time.Hour)

	return settingsMap, nil
}

func (uc *settingUseCase) Delete(key string) error {
	if err := uc.settingRepo.Delete(key); err != nil {
		return err
	}

	// Clear cache
	ctx := context.Background()
	_ = uc.cache.Delete(ctx, constants.SettingsCacheKey)

	return nil
}

func (uc *settingUseCase) updateSettingsCache() {
	ctx := context.Background()

	settings, err := uc.settingRepo.FindAll()
	if err == nil {
		settingsMap := make(map[string]string)
		for _, setting := range settings {
			settingsMap[setting.Key] = setting.Value
		}
		_ = uc.cache.Set(ctx, constants.SettingsCacheKey, settingsMap, 1*time.Hour)
	}
}
