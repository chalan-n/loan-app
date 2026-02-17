package models

type Religion struct {
	ReligionID   int    `json:"ReligionID" gorm:"column:ReligionID;primaryKey"`
	ReligionName string `json:"ReligionName" gorm:"column:ReligionName"`
}

func (Religion) TableName() string {
	return "religiontbl"
}
