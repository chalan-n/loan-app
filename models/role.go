// models/role.go
package models

// Role constants
const (
	RoleOfficer = "officer" // พนักงานสินเชื่อ — เห็นแค่งานตัวเอง
	RoleManager = "manager" // ผู้จัดการ — เห็นงานทุกคน + dashboard
	RoleAdmin   = "admin"   // ผู้ดูแลระบบ — จัดการ users + audit log
)
