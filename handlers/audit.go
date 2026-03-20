// handlers/audit.go
package handlers

import (
	"loan-app/config"
	"loan-app/models"

	"github.com/gofiber/fiber/v2"
)

// WriteAudit บันทึก audit log ทุกครั้งที่มีการกระทำสำคัญ
func WriteAudit(c *fiber.Ctx, action, refCode, detail string) {
	username := parseJWTUsername(c.Cookies("token"))
	role := getUserRole(username)
	ip := c.IP()
	ua := c.Get("User-Agent")
	if len(ua) > 255 {
		ua = ua[:255]
	}
	config.DB.Create(&models.AuditLog{
		Username:  username,
		Role:      role,
		Action:    action,
		RefCode:   refCode,
		Detail:    detail,
		IPAddress: ip,
		UserAgent: ua,
	})
}

// getUserRole ดึง role จาก DB ตาม username (cached ไม่ได้เพราะ role อาจเปลี่ยน)
func getUserRole(username string) string {
	if username == "" {
		return ""
	}
	var user models.User
	if err := config.DB.Select("role").Where("username = ?", username).First(&user).Error; err != nil {
		return models.RoleOfficer
	}
	return user.Role
}
