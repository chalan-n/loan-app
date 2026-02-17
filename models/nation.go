package models

type Nation struct {
	NationId   string `json:"NationId" gorm:"column:NationId;primaryKey"`
	NationName string `json:"NationName" gorm:"column:NationName"`
}

func (Nation) TableName() string {
	return "nation"
}
