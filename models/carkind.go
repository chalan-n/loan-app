package models

type CarKind struct {
	CarKindId   int    `json:"car_kind_id" gorm:"primaryKey;column:CarKindId"`
	CarKindName string `json:"car_kind_name" gorm:"column:CarKindName"`
}

func (CarKind) TableName() string {
	return "carkind"
}
