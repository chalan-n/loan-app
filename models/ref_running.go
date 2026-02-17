package models

type RefRunning struct {
	ID      int    `gorm:"primaryKey"`
	RefYear string `json:"ref_year"`
	EmpID   string `json:"emp_id"`
	Running int    `json:"running"`
}
