package models

type Tenant struct {
	BaseModel
	Name string `gorm:"type:varchar(100);not null" json:"name"`
}
