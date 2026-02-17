package models

type Showroom struct {
	ShowRoomId   string `json:"showroom_id" gorm:"primaryKey;column:ShowRoomId"`
	ShowRoomName string `json:"showroom_name" gorm:"column:ShowRoomName"`
}

func (Showroom) TableName() string {
	return "showroom"
}
