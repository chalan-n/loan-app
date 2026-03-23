# 📝 Loan Application System

ระบบจัดการคำขอสินเชื่อ พัฒนาด้วย Go (Golang) และ Fiber Framework

## 🛠️ เทคโนโลยีที่ใช้

- **Backend:** Go 1.23+ (Fiber Framework)
- **Database:** MySQL
- **ORM:** GORM
- **Template Engine:** HTML Templates
- **Authentication:** JWT + bcrypt

## 📋 ข้อกำหนดเบื้องต้น

1. **Go (Golang)** เวอร์ชัน 1.23 หรือสูงกว่า
   - ดาวน์โหลดได้ที่: https://golang.org/dl/

2. **MySQL Database**
   - ต้องมี MySQL Server ทำงานอยู่
   - สร้างฐานข้อมูลชื่อ `loan_db`

##📈 ฟีเจอร์ที่แนะนำเพิ่ม (เพื่อความสมบูรณ์)
1. ระบบสิทธิ์ผู้ใช้ (RBAC)
go
// models/role.go
type Role struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"unique"` // admin, officer, manager
    Permissions []string  // ["view_all", "approve", "delete"]
} 
// ผูกกับ User: RoleID uint
แบ่งสิทธิ์: Admin, Officer, Manager
จำกัดการเข้าถึงหน้า/ฟังก์ชันตาม Role

2. ระบบ Audit Log
go
type AuditLog struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint
    Action    string    // "create_loan", "approve", "delete"
    Detail    string    // JSON ของข้อมูลเก่า/ใหม่
    IPAddress string
    UserAgent string
    CreatedAt time.Time
}
บันทึกทุกการเปลี่ยนแปลงข้อมูล
ใช้สืบค้นปัญหา/ตรวจสอบ

3. ระบบ Notifications ขั้นสูง
แจ้งเตือนผ่าน LINE Notify / Email
แจ้งเตือนเมื่อใบคำขอใกล้หมดอายุ
แจ้งเตือนเมื่อมีใบคำขอใหม่

4. ระบบ Dashboard และรายงาน
Dashboard สำหรับ Manager (ยอดรวม/สถิติ/กราฟ)
รายงาน Excel ขั้นสูง (รายวัน/รายเดือน/ประเภทสินเชื่อ)
ระบบ Export PDF ใบคำขอ

5. ระบบ File Management
อัปโหลดไฟล์เอกสารเพิ่มเติม (สลิปเงินเดือน, บัญชีธนาคาร)
จัดเก็บไฟล์บน Cloud (S3/MinIO)
บีบอัดไฟล์อัตโนมัติ

6. ระบบ Workflow การอนุมัติ
หลายขั้นตอนการอนุมัติ (Officer → Manager → Director)
แจ้งเตือนเมื่อถึงคิวอนุมัติ
ประวัติการอนุมัติ

7. ระบบ Security ขั้นสูง
2FA (Google Authenticator)
Rate Limiting (ป้องกัน DoS)
Captcha บนหน้า Login
บันทึก IP ที่พยายามเข้าผิด

8. ระบบ Performance
แคชผลลัพธ์ OCR ด้วย Redis
ใช้งาน Database Connection Pool
บีบอัด Response ด้วย Gzip
CDN สำหรับ static files

9. ระบบ Mobile App (Flutter/React Native)
สแกนเอกสารผ่านกล้อง
แจ้งเตือน Push Notification
Offline Mode สำหรับพนักงานขาย

10. ระบบ Integration
ต่อกับระบบ Core Banking
ต่อกับระบบ Credit Bureau
API สำหรับระบบอื่นๆ เรียกใช้

------------------------------------------------------------------------------------------------------
🔧 การจัดการโครงสร้างให้สมบูรณ์
loan-app/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   ├── audit/
│   ├── notification/
│   └── workflow/
├── pkg/
│   ├── middleware/
│   └── utils/
├── configs/
├── migrations/
├── docs/
└── scripts/

------------------------------------------------------------------------------------------------------
📋 แผนพัฒนา (Roadmap)
Phase 1 
ระบบสิทธิ์ผู้ใช้ (RBAC)
ระบบ Audit Log
Dashboard สำหรับ Manager

Phase 2 
ระบบ File Management

Phase 3 
ระบบ Notifications ขั้นสูง (LINE/Email)
ระบบ Mobile App
ระบบ Integration 
ระบบ Reporting ขั้นสูง
