// handlers/webauthn.go
package handlers

import (
	"encoding/base64"
	"encoding/json"
	"loan-app/config"
	"loan-app/models"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
)

// ── WebAuthn instance (init ใน InitWebAuthn) ────────────────────────────────
var WebAuthn *webauthn.WebAuthn

func InitWebAuthn(rpID, rpOrigin, rpName string) error {
	var err error
	WebAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: rpName,
		RPID:          rpID,
		RPOrigins:     []string{rpOrigin},
	})
	return err
}

// ── webAuthnUser implements webauthn.User ────────────────────────────────────
type webAuthnUser struct {
	user  models.User
	creds []webauthn.Credential
}

func (u *webAuthnUser) WebAuthnID() []byte {
	id := make([]byte, 8)
	uid := uint64(u.user.ID)
	for i := 0; i < 8; i++ {
		id[7-i] = byte(uid >> (i * 8))
	}
	return id
}
func (u *webAuthnUser) WebAuthnName() string        { return u.user.Username }
func (u *webAuthnUser) WebAuthnDisplayName() string { return u.user.Username }
func (u *webAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.creds
}

// loadWebAuthnUser โหลด user และ credentials จาก DB
func loadWebAuthnUser(username string) (*webAuthnUser, error) {
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	var dbCreds []models.WebAuthnCredential
	config.DB.Where("user_id = ?", user.ID).Find(&dbCreds)

	creds := make([]webauthn.Credential, 0, len(dbCreds))
	for _, dc := range dbCreds {
		credIDBytes, err := base64.RawURLEncoding.DecodeString(dc.CredentialID)
		if err != nil {
			continue
		}

		var transports []protocol.AuthenticatorTransport
		if dc.Transport != "" {
			for _, t := range strings.Split(dc.Transport, ",") {
				transports = append(transports, protocol.AuthenticatorTransport(strings.TrimSpace(t)))
			}
		}

		creds = append(creds, webauthn.Credential{
			ID:              credIDBytes,
			PublicKey:       dc.PublicKey,
			AttestationType: dc.AttestationType,
			Transport:       transports,
			Authenticator: webauthn.Authenticator{
				SignCount: dc.SignCount,
			},
		})
	}

	return &webAuthnUser{user: user, creds: creds}, nil
}

// webAuthnReady ตรวจสอบว่า WebAuthn instance พร้อมใช้งาน
func webAuthnReady(c *fiber.Ctx) bool {
	if WebAuthn == nil {
		c.Status(503).JSON(fiber.Map{"error": "WebAuthn ยังไม่พร้อม กรุณาตั้งค่า WEBAUTHN_RPID และ WEBAUTHN_ORIGIN ใน .env แล้ว restart"})
		return false
	}
	return true
}

// ── Register: Begin ──────────────────────────────────────────────────────────
// POST /webauthn/register/begin
func WebAuthnRegisterBegin(c *fiber.Ctx) error {
	if !webAuthnReady(c) {
		return nil
	}
	username := parseJWTUsername(c.Cookies("token"))
	if username == "" {
		return c.Status(401).JSON(fiber.Map{"error": "กรุณาเข้าสู่ระบบก่อน"})
	}

	waUser, err := loadWebAuthnUser(username)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "ไม่พบผู้ใช้"})
	}

	options, sessionData, err := WebAuthn.BeginRegistration(waUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "เริ่มลงทะเบียนไม่สำเร็จ: " + err.Error()})
	}

	// เก็บ sessionData ใน cookie (JSON base64)
	sdJSON, _ := json.Marshal(sessionData)
	c.Cookie(&fiber.Cookie{
		Name:     "wa_reg_session",
		Value:    base64.StdEncoding.EncodeToString(sdJSON),
		HTTPOnly: true,
		Path:     "/",
		MaxAge:   300, // 5 นาที
	})

	return c.JSON(options)
}

// ── Register: Finish ─────────────────────────────────────────────────────────
// POST /webauthn/register/finish
func WebAuthnRegisterFinish(c *fiber.Ctx) error {
	if !webAuthnReady(c) {
		return nil
	}
	username := parseJWTUsername(c.Cookies("token"))
	if username == "" {
		return c.Status(401).JSON(fiber.Map{"error": "กรุณาเข้าสู่ระบบก่อน"})
	}

	// กู้คืน sessionData จาก cookie
	sessionCookie := c.Cookies("wa_reg_session")
	if sessionCookie == "" {
		return c.Status(400).JSON(fiber.Map{"error": "session หมดอายุ กรุณาลองใหม่"})
	}
	sdJSON, err := base64.StdEncoding.DecodeString(sessionCookie)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "session ไม่ถูกต้อง"})
	}
	var sessionData webauthn.SessionData
	if err := json.Unmarshal(sdJSON, &sessionData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "session ไม่ถูกต้อง"})
	}

	waUser, err := loadWebAuthnUser(username)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "ไม่พบผู้ใช้"})
	}

	// Parse body เป็น io.Reader สำหรับ webauthn library
	bodyReader := strings.NewReader(string(c.Body()))
	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(bodyReader)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ข้อมูล credential ไม่ถูกต้อง: " + err.Error()})
	}

	credential, err := WebAuthn.CreateCredential(waUser, sessionData, parsedResponse)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ลงทะเบียนไม่สำเร็จ: " + err.Error()})
	}

	// แปลง transport เป็น string
	transportStrs := make([]string, len(credential.Transport))
	for i, t := range credential.Transport {
		transportStrs[i] = string(t)
	}

	// บันทึกลง DB
	dbCred := models.WebAuthnCredential{
		UserID:          waUser.user.ID,
		CredentialID:    base64.RawURLEncoding.EncodeToString(credential.ID),
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		Transport:       strings.Join(transportStrs, ","),
		SignCount:       credential.Authenticator.SignCount,
		DeviceName:      c.Get("User-Agent"),
	}
	if err := config.DB.Create(&dbCred).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "บันทึก credential ไม่สำเร็จ"})
	}

	// ลบ session cookie
	c.ClearCookie("wa_reg_session")

	WriteAudit(c, "webauthn_register", "", "ลงทะเบียน biometric สำเร็จ")
	return c.JSON(fiber.Map{"success": true, "message": "ลงทะเบียนสำเร็จ"})
}

// ── Login: Begin ─────────────────────────────────────────────────────────────
// POST /webauthn/login/begin  body: {"username": "570639"}
func WebAuthnLoginBegin(c *fiber.Ctx) error {
	if !webAuthnReady(c) {
		return nil
	}
	type Req struct {
		Username string `json:"username"`
	}
	var req Req
	if err := c.BodyParser(&req); err != nil || req.Username == "" {
		return c.Status(400).JSON(fiber.Map{"error": "กรุณาระบุ username"})
	}

	waUser, err := loadWebAuthnUser(req.Username)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "ไม่พบผู้ใช้"})
	}
	if len(waUser.creds) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ยังไม่ได้ลงทะเบียน biometric"})
	}

	options, sessionData, err := WebAuthn.BeginLogin(waUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "เริ่ม login ไม่สำเร็จ: " + err.Error()})
	}

	// เก็บ sessionData + username ใน cookie
	type loginSession struct {
		Username    string               `json:"username"`
		SessionData webauthn.SessionData `json:"session_data"`
	}
	ls := loginSession{Username: req.Username, SessionData: *sessionData}
	lsJSON, _ := json.Marshal(ls)
	c.Cookie(&fiber.Cookie{
		Name:     "wa_login_session",
		Value:    base64.StdEncoding.EncodeToString(lsJSON),
		HTTPOnly: true,
		Path:     "/",
		MaxAge:   300,
	})

	return c.JSON(options)
}

// ── Login: Finish ─────────────────────────────────────────────────────────────
// POST /webauthn/login/finish
func WebAuthnLoginFinish(c *fiber.Ctx) error {
	if !webAuthnReady(c) {
		return nil
	}
	// กู้คืน session
	sessionCookie := c.Cookies("wa_login_session")
	if sessionCookie == "" {
		return c.Status(400).JSON(fiber.Map{"error": "session หมดอายุ กรุณาลองใหม่"})
	}
	lsJSON, err := base64.StdEncoding.DecodeString(sessionCookie)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "session ไม่ถูกต้อง"})
	}

	type loginSession struct {
		Username    string               `json:"username"`
		SessionData webauthn.SessionData `json:"session_data"`
	}
	var ls loginSession
	if err := json.Unmarshal(lsJSON, &ls); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "session ไม่ถูกต้อง"})
	}

	waUser, err := loadWebAuthnUser(ls.Username)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "ไม่พบผู้ใช้"})
	}

	bodyReader := strings.NewReader(string(c.Body()))
	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(bodyReader)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ข้อมูลไม่ถูกต้อง: " + err.Error()})
	}

	credential, err := WebAuthn.ValidateLogin(waUser, ls.SessionData, parsedResponse)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "ยืนยันตัวตนไม่สำเร็จ: " + err.Error()})
	}

	// อัปเดต SignCount
	credIDStr := base64.RawURLEncoding.EncodeToString(credential.ID)
	config.DB.Model(&models.WebAuthnCredential{}).
		Where("user_id = ? AND credential_id = ?", waUser.user.ID, credIDStr).
		Update("sign_count", credential.Authenticator.SignCount)

	// สร้าง JWT เหมือน LoginPost
	c.ClearCookie("wa_login_session")
	tokenStr, err := createJWTToken(ls.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "สร้าง token ไม่สำเร็จ"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenStr,
		HTTPOnly: true,
		Path:     "/",
	})

	WriteAuditAs(c, ls.Username, "login_biometric", "", "เข้าสู่ระบบด้วย biometric สำเร็จ")
	return c.JSON(fiber.Map{"success": true, "redirect": "/main"})
}

// ── Register Page ────────────────────────────────────────────────────────────
// GET /webauthn/register
func WebAuthnRegisterPage(c *fiber.Ctx) error {
	username := parseJWTUsername(c.Cookies("token"))
	if username == "" {
		return c.Redirect("/login")
	}
	var creds []models.WebAuthnCredential
	var user models.User
	config.DB.Select("id").Where("username = ?", username).First(&user)
	config.DB.Where("user_id = ?", user.ID).Find(&creds)

	return c.Render("webauthn_register", fiber.Map{
		"Username":    username,
		"Credentials": creds,
		"CurrentRole": GetCurrentUserRole(c),
	})
}

// ── List & Delete credentials ─────────────────────────────────────────────────
// GET /webauthn/credentials
func WebAuthnListCredentials(c *fiber.Ctx) error {
	username := parseJWTUsername(c.Cookies("token"))
	if username == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var user models.User
	config.DB.Select("id").Where("username = ?", username).First(&user)

	var creds []models.WebAuthnCredential
	config.DB.Where("user_id = ?", user.ID).Find(&creds)
	return c.JSON(creds)
}

// DELETE /webauthn/credentials/:id
func WebAuthnDeleteCredential(c *fiber.Ctx) error {
	username := parseJWTUsername(c.Cookies("token"))
	if username == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var user models.User
	config.DB.Select("id").Where("username = ?", username).First(&user)

	credID := c.Params("id")
	if err := config.DB.Where("id = ? AND user_id = ?", credID, user.ID).
		Delete(&models.WebAuthnCredential{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "ลบไม่สำเร็จ"})
	}
	WriteAudit(c, "webauthn_delete_credential", "", "credentialID="+credID)
	return c.JSON(fiber.Map{"success": true})
}
