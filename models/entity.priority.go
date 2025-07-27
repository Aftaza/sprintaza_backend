package model

import (
	"time"

	"gorm.io/gorm"
)

type EntityPriority struct {
	ID       	uint 			`gorm:"primaryKey;column:priority_id"` // PK
	ProjectID	uint			`gorm:"column:project_id;not null"` // FK ke Project table
	Title     	string			`gorm:"column:title;type:varchar(100);not null"`

	CreatedAt 	time.Time      `gorm:"column:created_at"`
	UpdatedAt 	time.Time      `gorm:"column:updated_at"`
	DeletedAt 	gorm.DeletedAt `gorm:"index;column:deleted_at"` // Opsi: untuk soft delete

	// Relasi GORM
	Project 	EntityProject	`gorm:"foreignKey:ProjectID"`
}

// TableName mengembalikan nama nama tabel di database
func (EntityPriority) TableName() string {
	return "PRIORITY"
}