package model

import (
	"time"

	"gorm.io/gorm"
)

type EntitySubtask struct {
	ID       uint 		`gorm:"primaryKey;column:subtask_id"` // PK
	TaskID 	 uint 		`gorm:"column:task_id;not null"` 		// FK ke task
	Title	 string 	`gorm:"column:title;type:varchar(100);not null"`

	CreatedAt 	time.Time      `gorm:"column:created_at"`
	UpdatedAt 	time.Time      `gorm:"column:updated_at"`
	DeletedAt 	gorm.DeletedAt `gorm:"index;column:deleted_at"` // Opsi: untuk soft delete

	// Relasi
	Task EntityTask `gorm:"foreignKey:TaskID"`
}

// TableName mengembalikan nama nama tabel di database
func (EntitySubtask) TableName() string {
	return "SUBTASK"
}