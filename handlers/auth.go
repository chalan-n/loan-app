package handlers

import (
	"loan-app/config"
	"loan-app/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginPage(c *fiber.Ctx) error {
	// Check if user is already logged in
	if parseJWTUsername(c.Cookies("token")) != "" {
		return c.Redirect("/main")
	}

	// ส่ง error message ไปหน้า login ถ้ามี (จาก login ไม่สำเร็จ)
	errorMsg := c.Query("error")
	return c.Render("login", fiber.Map{
		"Error": errorMsg,
	})
}

func LoginPost(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user := new(models.User)
	config.DB.Where("username = ?", username).First(user)

	if user.ID == 0 || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		// ล็อกอินผิด → ส่งกลับไป login พร้อมสถานะ error
		return c.Redirect("/login?status=error")
	}

	// ล็อกอินสำเร็จ — สร้าง session ID ใหม่ (kick session เก่า)
	newSessionID := uuid.NewString()
	config.DB.Model(user).Update("current_session_id", newSessionID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   user.Username,
		"session_id": newSessionID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(config.GetConfig().JWTSecret))

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenStr,
		HTTPOnly: true,
		Path:     "/",
	})

	// ส่งกลับไป login พร้อมสถานะ success → JS จะโชว์ Toast เขียว แล้วค่อยไป /main
	return c.Redirect("/main")
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("token")
	return c.Redirect("/login")
}

func AuthMiddleware(c *fiber.Ctx) error {
	tokenStr := c.Cookies("token")
	if tokenStr == "" {
		return c.Redirect("/login")
	}
	username := parseJWTUsername(tokenStr)
	if username == "" {
		c.ClearCookie("token")
		return c.Redirect("/login")
	}

	// ตรวจสอบ session_id: ถ้าไม่ตรงกับ DB แปลว่ามีการล็อกอินจากอุปกรณ์ใหม่
	claimSessionID := parseJWTSessionID(tokenStr)
	if claimSessionID != "" {
		var user models.User
		if err := config.DB.Select("current_session_id").Where("username = ?", username).First(&user).Error; err == nil {
			if user.CurrentSessionID != claimSessionID {
				c.ClearCookie("token")
				return c.Redirect("/login?error=session_expired")
			}
		}
	}

	return c.Next()
}

func ChangePasswordPage(c *fiber.Ctx) error {
	return c.Render("change_password", nil)
}

// MobileAPIKeyMiddleware ตรวจ X-API-Key header สำหรับ endpoint ที่ mobile app เรียก
func MobileAPIKeyMiddleware(c *fiber.Ctx) error {
	key := c.Get("X-API-Key")
	expected := config.GetConfig().MobileAPIKey
	if expected == "" || key != expected {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: invalid or missing API key"})
	}
	return c.Next()
}

func ChangePasswordPost(c *fiber.Ctx) error {
	oldPassword := c.FormValue("old_password")
	newPassword := c.FormValue("new_password")

	// Get user from token
	username := parseJWTUsername(c.Cookies("token"))

	if username == "" {
		// Should be caught by middleware normally, but double check
		return c.Status(401).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"success": false, "message": "User not found"})
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return c.JSON(fiber.Map{"success": false, "message": "รหัสผ่านเดิมไม่ถูกต้อง"})
	}

	// Hash new password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	user.Password = string(hashedPassword)

	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(fiber.Map{"success": false, "message": "บันทึกข้อมูลไม่สำเร็จ"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "เปลี่ยนรหัสผ่านเรียบร้อยแล้ว"})
}
