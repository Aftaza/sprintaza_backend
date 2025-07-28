package model

import (
	"time"

	"gorm.io/gorm"
)

type EntityColumn struct {
	ID       	uint 	`gorm:"primaryKey;column:column_id"` // PK
	Name     	string 	`gorm:"column:name;type:varchar(100);not null"`
	ProjectID	uint	`gorm:"column:task_id;not null"` 		// FK ke project
	Order 		int		`gorm:"column:order;not null"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"` // Opsi: untuk soft delete

	// Relasi
	Project EntityProject `gorm:"foreignKey:ProjectID"`
}

// TableName mengembalikan nama nama tabel di database
func (EntityColumn) TableName() string {
	return "columns"
}

// DefaultColumnNames adalah slice dari nama-nama kolom default
var DefaultColumnNames = []string{
	"To Do",
	"In Progress",
	"Done",
}