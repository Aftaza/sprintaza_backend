package model

// User merepresentasikan entitas pengguna dalam sistem.
type User struct {
	ID                 uint              `gorm:"primaryKey" json:"id"`
	Name               string            `gorm:"type:varchar(100);not null" json:"name"`
	Email              string            `gorm:"type:varchar(100);unique;not null" json:"email"`
	PasswordHash       string            `gorm:"type:varchar(255);not null" json:"-"`
	AuthProvider       string            `gorm:"type:varchar(50)" json:"auth_provider,omitempty"`
	UserXPID           uint              `json:"-"`
	UserXP             UserXP            `gorm:"foreignKey:UserXPID" json:"xp_details"`
	ProjectMemberships []ProjectMember   `gorm:"foreignKey:UserID" json:"project_memberships,omitempty"`
	UserAchievements   []UserAchievement `gorm:"foreignKey:UserID" json:"user_achievements,omitempty"`
	AssignedTasks      []Task            `gorm:"foreignKey:AssignedUserID" json:"assigned_tasks,omitempty"`
}

// UserXP menyimpan data gamifikasi poin untuk setiap pengguna.
type UserXP struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	UserID    uint `json:"user_id"`
	XP        int  `gorm:"default:0" json:"xp"`
	InitialXP int  `gorm:"default:0" json:"initial_xp"`
}