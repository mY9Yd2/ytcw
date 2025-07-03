package schema

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Channels []Channel `gorm:"foreignKey:CategoryID;"`

	Name string `gorm:"unique;size:40;"`
}
