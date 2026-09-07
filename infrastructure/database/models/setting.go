package models

type SettingModel struct {
	Key   string `gorm:"primaryKey;index"`
	Value string `json:"value"`
}
