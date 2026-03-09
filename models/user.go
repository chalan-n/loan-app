// models/user.go
package models

type User struct {
	ID               uint   `gorm:"primaryKey"`
	Username         string `gorm:"unique"`
	Password         string
	CurrentSessionID string `gorm:"type:varchar(36)"`
}
