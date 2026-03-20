// models/user.go
package models

type User struct {
	ID               uint   `gorm:"primaryKey"`
	Username         string `gorm:"unique"`
	Password         string
	Role             string `gorm:"type:varchar(20);default:'officer'"`
	CurrentSessionID string `gorm:"type:varchar(36)"`
}
