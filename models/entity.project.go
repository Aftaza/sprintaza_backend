package model

import (
	"time"

	"gorm.io/gorm"
)

type EntityProject struct {
	ID       	uint 			`gorm:"primaryKey;column:project_id"` // PK
	Name     	string 			`gorm:"column:name;type:varchar(100);not null"`
	Description string 			`gorm:"column:description;type:text"`

	CreatedAt 	time.Time      `gorm:"column:created_at"`
	UpdatedAt 	time.Time      `gorm:"column:updated_at"`
	DeletedAt 	gorm.DeletedAt `gorm:"index;column:deleted_at"` // Opsi: untuk soft delete
}

// TableName mengembalikan nama nama tabel di database
func (EntityProject) TableName() string {
	return "PROJECTS"
}