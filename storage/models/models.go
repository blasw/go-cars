package models

type Cars struct {
	ID      int    `gorm:"primaryKey"`
	RegNum  string `gorm:"unique"`
	Mark    string `gorm:"not null"`
	Model   string `gorm:"not null"`
	Year    int    `gorm:"not null"`
	OwnerID int    `gorm:"not null"`
	Owner   Owners `gorm:"foreignKey:OwnerID"`
}

type Owners struct {
	ID         int    `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	Surname    string `gorm:"not null"`
	Patronymic string `gorm:"not null"`
}
