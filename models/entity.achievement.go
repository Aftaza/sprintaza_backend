package model

import (
	"time"

	"gorm.io/gorm"
)

// Achievement merepresentasikan tabel ACHIEVEMENTS
type EntityAchievement struct {
	ID          uint   `gorm:"primaryKey;column:achievement_id"` // PK
	Name        string `gorm:"column:name;type:varchar(100);not null"`
	Description string `gorm:"column:description;type:text"`
	XPReward    int    `gorm:"column:xp_reward;not null"`
	IconURL     string `gorm:"column:icon_url;type:varchar(255)"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"` // Opsi: untuk soft delete
}

// TableName mengembalikan nama nama tabel di database
func (EntityAchievement) TableName() string {
	return "achievements"
}

// DefaultAchievements adalah slice dari achievement default yang akan di-seed
var DefaultAchievements = []EntityAchievement{
	{Name: "Selamat Datang!", Description: "Diberikan saat kamu pertama kali bergabung. Awal dari perjalanan produktifmu!", XPReward: 10, IconURL: "rocket"},
	{Name: "Perencana Proyek", Description: "Kamu berhasil membuat proyek pertamamu. Saatnya mengubah ide menjadi aksi.", XPReward: 25, IconURL: "notebook-pen"},
	{Name: "Tugas Pertama", Description: "Menyelesaikan tugas pertamamu. Satu langkah kecil yang memulai kesuksesan besar.", XPReward: 15, IconURL: "folder-check"},
	{Name: "Anggota Tim", Description: "Berhasil mengundang anggota tim pertamamu ke dalam sebuah proyek. Kolaborasi adalah kunci!", XPReward: 30, IconURL: "handshake"},
	{Name: "Pemanasan", Description: "Menyelesaikan 5 tugas. Mesin produktivitasmu mulai panas!", XPReward: 30, IconURL: "flame"},
	{Name: "Koneksi Kode", Description: "Menghubungkan proyek dengan repositori Git (misalnya GitHub/GitLab) untuk pertama kalinya.", XPReward: 50, IconURL: "link"},
	{Name: "Delegator Andal", Description: "Tugas yang kamu delegasikan berhasil diselesaikan oleh rekan timmu. Kepercayaanmu terbayar!", XPReward: 20, IconURL: "shield-check"},
	{Name: "Got Your Back!", Description: "Menyelesaikan tugas yang didelegasikan kepadamu. Kamu adalah rekan tim yang bisa diandalkan!", XPReward: 35, IconURL: "sticker"},
	{Name: "Kontributor", Description: "Berhasil menyelesaikan 25 tugas. Kamu adalah kontributor kunci bagi kesuksesan tim!", XPReward: 80, IconURL: "gem"},
}