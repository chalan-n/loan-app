package models

type Title struct {
	TitleID   string `json:"TitleID" gorm:"column:TitleID;primaryKey"`
	TitleName string `json:"TitleName" gorm:"column:TitleName"`
	TitleType string `json:"TitleType" gorm:"column:TitleType"`
}

func (Title) TableName() string {
	return "title"
}
