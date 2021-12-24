package db

type PhoneNumber struct {
	ID    uint `gorm:"primaryKey"`
	Value string `gorm:"unique"`
}
