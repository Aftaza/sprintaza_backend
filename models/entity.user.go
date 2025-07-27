package model

import (
	"time"

	util "github.com/Aftaza/sprintaza_backend/utils"
	"gorm.io/gorm"
)

type EntityUsers struct {
	ID 				uint	`gorm:"primaryKey;column:user_id"` // PK
	Name			string	`gorm:"column:name;type:varchar(100);not null"`
	Email         	string 	`gorm:"column:email;type:varchar(100);not null"`
	PasswordHash  	string 	`gorm:"column:password_hash;type:varchar(255);not null"`
	AvatarURL     	string 	`gorm:"column:avatar_url;type:varchar(100)"` // AvatarURL bisa null
	UserXPID      	uint	`gorm:"column:user_xp_id"`     // FK ke UserXP

	CreatedAt time.Time      `gorm:"column:created_at"` // Otomatis dikelola oleh GORM
	UpdatedAt time.Time      `gorm:"column:updated_at"` // Otomatis dikelola oleh GORM
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"` // Opsi: untuk soft delete dengan GORM

	// Relasi GORM
	UserXP    EntityUserXP	 `gorm:"foreignKey:UserXPID"` // Relasi One-to-One dengan UserXP
}

func (EntityUsers) TableName() string {
	return "USERS"
}

func (entity *EntityUsers) BeforeCreate(db *gorm.DB) (err error) {
	entity.PasswordHash = util.HashPassword(entity.PasswordHash)
	return nil
}