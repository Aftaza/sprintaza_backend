package model

import (
	"time"

	"gorm.io/gorm"
)

// UserAchievement merepresentasikan tabel USER_ACHIEVEMENTS (tabel penghubung Many-to-Many)
type EntityUserAchievement struct {
	ID 				  uint `gorm:"primaryKey;column:user_achievement_id"` // PK
	UserID            uint `gorm:"column:user_id;not null"`           // FK ke USERS
	AchievementID     uint `gorm:"column:achievement_id;not null"`    // FK ke ACHIEVEMENTS
	UnlockedAt        time.Time `gorm:"column:unlocked_at;not null"` // Waktu achievement didapat

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"` // Opsi: untuk soft delete

	// Relasi GORM
	User        EntityUsers        `gorm:"foreignKey:UserID"`        // Relasi Many-to-One dengan User
	Achievement EntityAchievement `gorm:"foreignKey:AchievementID"` // Relasi Many-to-One dengan Achievement
}

// TableName mengembalikan nama tabel di database
func (EntityUserAchievement) TableName() string {
	return "USER_ACHIEVEMENTS"
}