// config/database.go
package config

import (
	"loan-app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:@tcp(127.0.0.1:3308)/loan_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("เชื่อม DB ไม่ได้: " + err.Error())
	}
	DB = db
	db.AutoMigrate(&models.User{}, &models.LoanApplication{}, &models.Guarantor{})
}
