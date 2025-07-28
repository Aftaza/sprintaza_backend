package model

import (
	"time"

	"gorm.io/gorm"
)

// UserXP merepresentasikan tabel USER_XP
type EntityUserXP struct {
	ID       uint `gorm:"primaryKey;column:user_xp_id"` // PK
	TotalXP  int  `gorm:"column:total_xp;not null"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"` // Opsi: untuk soft delete
	// Relasi GORM (optional, bisa dihandle dari sisi User)
	// User User `gorm:"foreignKey:UserXPID"` // Relasi One-to-One kembali ke User
}

// TableName mengembalikan nama tabel di database
func (EntityUserXP) TableName() string {
	return "user_xp"
}