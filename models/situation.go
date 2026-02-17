package models

type Situation struct {
	SituationID   int    `json:"SituationID" gorm:"column:SituationID;primaryKey"`
	SituationName string `json:"SituationName" gorm:"column:SituationName"`
}

func (Situation) TableName() string {
	return "situation"
}
