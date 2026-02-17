package models

type Occupy struct {
	OccupyID   int    `json:"OccupyID" gorm:"column:OccupyID;primaryKey"`
	OccupyName string `json:"OccupyName" gorm:"column:OccupyName"`
}

func (Occupy) TableName() string {
	return "occupy"
}
