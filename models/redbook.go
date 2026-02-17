package models

type RedbookII struct {
	ID          uint    `gorm:"primaryKey"`
	CarCode     string  `json:"car_code" gorm:"column:carCode;uniqueIndex"`
	CarType     string  `json:"car_type" gorm:"column:carType"`
	CarBrand    string  `json:"car_brand" gorm:"column:carBrand"`
	CarSubModel string  `json:"car_sub_model" gorm:"column:carSubModel"`
	CarGear     string  `json:"car_gear" gorm:"column:carGear"`
	CarYear     string  `json:"car_year" gorm:"column:carYear"`
	CarPrice    float64 `json:"car_price" gorm:"column:carPrice"`
}

// TableName overrides the table name used by User to `redbook_ii`
func (RedbookII) TableName() string {
	return "redbook_ii"
}
