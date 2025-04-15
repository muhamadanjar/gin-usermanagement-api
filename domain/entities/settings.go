package entities

type Setting struct {
	Key   string `gorm:"primaryKey;index"`
	Value string `json:"value"`
}
