package models

type Tenant struct {
	BaseModel
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	VillageCode string `gorm:"type:varchar(20);not null;unique" json:"village_code"`
}
