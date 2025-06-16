package model

import "time"

// Achievement merepresentasikan daftar pencapaian yang bisa didapat.
type Achievement struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	XpReward    int    `gorm:"not null" json:"xp_reward"`
	IconURL     string `gorm:"type:varchar(255)" json:"icon_url"`
}

// UserAchievement menandakan seorang user telah membuka sebuah achievement.
type UserAchievement struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	UserID        uint        `gorm:"not null" json:"user_id"`
	AchievementID uint        `gorm:"not null" json:"achievement_id"`
	UnlockedAt    time.Time   `gorm:"not null" json:"unlocked_at"`
	User          User        `gorm:"foreignKey:UserID" json:"-"`
	Achievement   Achievement `gorm:"foreignKey:AchievementID" json:"achievement"`
}