package models

type InsuClass struct {
	InsuranceClassID   string `json:"insurance_classid" gorm:"column:insurance_classid;primaryKey"`
	InsuranceClassName string `json:"insurance_classname" gorm:"column:insurance_classname"`
}

func (InsuClass) TableName() string {
	return "insuclss"
}
