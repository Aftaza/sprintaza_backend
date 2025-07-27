package model

import (
	"time"

	"gorm.io/gorm"
)

type EntityProjectMember struct {
	UserID      uint `gorm:"primaryKey;column:user_id;"` // PK
	ProjectID   uint `gorm:"primaryKey;column:project_id;"` // PK
	RoleID 		uint `gorm:"column:role_id;not null"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"` // Opsi: untuk soft delete

	// relasi
	User 	EntityUsers 	`gorm:"foreignKey:UserID"`
	Project EntityProject 	`gorm:"foreignKey:ProjectID"`
	Role 	EntityRole 		`gorm:"foreignKey:RoleID"`
}

// TableName mengembalikan nama nama tabel di database
func (EntityProjectMember) TableName() string {
	return "PROJECT_MEMBER"
}