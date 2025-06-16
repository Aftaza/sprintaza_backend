package model

import "time"

// Project merepresentasikan sebuah proyek.
type Project struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	Name        string          `gorm:"type:varchar(100);not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	Members     []ProjectMember `gorm:"foreignKey:ProjectID" json:"members,omitempty"`
	Columns     []Column        `gorm:"foreignKey:ProjectID" json:"columns,omitempty"`
}

// Role merepresentasikan peran pengguna dalam sebuah proyek.
type Role struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description string `gorm:"type:text" json:"description"`
}

// ProjectMember adalah tabel penghubung antara User, Project, dan Role.
type ProjectMember struct {
	UserID    uint `gorm:"primaryKey" json:"user_id"`
	ProjectID uint `gorm:"primaryKey" json:"project_id"`
	RoleID    uint `gorm:"not null" json:"role_id"`
	User      User `gorm:"foreignKey:UserID" json:"user"`
	Project   Project `gorm:"foreignKey:ProjectID" json:"project"`
	Role      Role `gorm:"foreignKey:RoleID" json:"role"`
}

// Column merepresentasikan sebuah kolom di dalam board proyek.
type Column struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProjectID uint   `gorm:"not null" json:"project_id"`
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	Order     int    `gorm:"not null" json:"order"`
	Tasks     []Task `gorm:"foreignKey:ColumnID" json:"tasks,omitempty"`
}

// Task merepresentasikan sebuah tugas di dalam sebuah Column.
type Task struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ColumnID       uint       `gorm:"not null" json:"column_id"`
	Title          string     `gorm:"type:varchar(100);not null" json:"title"`
	Description    string     `gorm:"type:text" json:"description"`
	AssignedUserID *uint      `json:"assigned_user_id"`
	StartDateTime  *time.Time `json:"start_datetime,omitempty"`
	EndDateTime    *time.Time `json:"end_datetime,omitempty"`
	AssignedUser   *User      `gorm:"foreignKey:AssignedUserID" json:"assigned_user,omitempty"`
}