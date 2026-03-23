package models

// WebAuthnCredential เก็บ public key credential ของ user แต่ละอุปกรณ์
type WebAuthnCredential struct {
	ID              uint   `gorm:"primaryKey"`
	UserID          uint   `gorm:"not null;index"`
	CredentialID    string `gorm:"type:varchar(512);not null;uniqueIndex"`
	PublicKey       []byte `gorm:"type:blob;not null"`
	AttestationType string `gorm:"type:varchar(64)"`
	Transport       string `gorm:"type:varchar(128)"`
	SignCount       uint32 `gorm:"default:0"`
	DeviceName      string `gorm:"type:varchar(512)"` // ชื่ออุปกรณ์ (user-agent)
	CreatedAt       LocalTime
}
