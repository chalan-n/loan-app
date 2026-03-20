// models/audit_log.go
package models

import "time"

type AuditLog struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"type:varchar(50);index"`
	Role      string    `gorm:"type:varchar(20)"`
	Action    string    `gorm:"type:varchar(100)"` // e.g. "create_loan", "delete_loan", "login"
	RefCode   string    `gorm:"type:varchar(50)"`  // loan ref code ถ้ามี
	Detail    string    `gorm:"type:text"`         // รายละเอียดเพิ่มเติม
	IPAddress string    `gorm:"type:varchar(45)"`
	UserAgent string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `gorm:"index"`
}
