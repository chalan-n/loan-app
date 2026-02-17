package models

type Race struct {
	RaceID   int    `json:"RaceID" gorm:"column:RaceID;primaryKey"`
	RaceName string `json:"RaceName" gorm:"column:RaceName"`
}

func (Race) TableName() string {
	return "race"
}
