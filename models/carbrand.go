package models

type CarBrand struct {
	CarBrandID   int    `json:"car_brand_id" gorm:"primaryKey;column:CarBrandID"`
	CarBrandName string `json:"car_brand_name" gorm:"column:CarBrandName"`
}

func (CarBrand) TableName() string {
	return "carbrand"
}
