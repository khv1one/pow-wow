package entity

type Quote struct {
	ID     int64   `gorm:"primary_key;auto_increment"`
	Quote  string  `gorm:"type:text"`
	Author *string `gorm:"type:VARCHAR(255)"`
}
