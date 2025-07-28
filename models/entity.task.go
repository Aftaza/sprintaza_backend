package model

import (
	"time"

	"gorm.io/gorm"
)

type EntityTask struct {
	ID       	uint		`gorm:"primaryKey;column:task_id"` // PK
	Title     	string		`gorm:"column:title;type:varchar(100);not null"`
	Description string		`gorm:"column:description;type:text"`
	Priority	uint		`gorm:"column:priority_id;not null"`           		// FK ke priority
	Assignee	uint 		`gorm:"column:assignee_user_id;not null"`           // FK ke USERS
	StartTime	time.Time	`gorm:"column:start_datetime;"`
	EndTime		time.Time	`gorm:"column:end_datetime;not null"`

	CreatedAt 	time.Time      `gorm:"column:created_at"`
	UpdatedAt 	time.Time      `gorm:"column:updated_at"`
	DeletedAt 	gorm.DeletedAt `gorm:"index;column:deleted_at"` // Opsi: untuk soft delete
}

// TableName mengembalikan nama nama tabel di database
func (EntityTask) TableName() string {
	return "task"
}