package models

type InsuComp struct {
	CompanyID   int    `json:"CompanyID" gorm:"column:CompanyID;primaryKey"`
	CompanyName string `json:"CompanyName" gorm:"column:CompanyName"`
}

func (InsuComp) TableName() string {
	return "insucomp"
}
