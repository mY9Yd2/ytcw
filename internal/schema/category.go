package schema

type Category struct {
	UUIDModel
	Channels []Channel `gorm:"foreignKey:CategoryRefer;"`

	Name string `gorm:"unique;size:40;"`
}
