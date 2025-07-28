package model

import (
	"time"

	"gorm.io/gorm"
)

type EntityRole struct {
	ID       	uint 	`gorm:"primaryKey;column:role_id"` // PK
	Name     	string 	`gorm:"column:name;type:varchar(100);not null"`
	Description string 	`gorm:"column:description;type:text"`

	CreatedAt 	time.Time      `gorm:"column:created_at"`
	UpdatedAt 	time.Time      `gorm:"column:updated_at"`
	DeletedAt 	gorm.DeletedAt `gorm:"index;column:deleted_at"` // Opsi: untuk soft delete
}

// TableName mengembalikan nama nama tabel di database
func (EntityRole) TableName() string {
	return "roles"
}

var DefaultRoles = []EntityRole{
	{Name: "Owner", Description: "Kontrol penuh dan otoritas tertinggi atas seluruh sistem/proyek. Bertanggung jawab atas arah strategis, kepemilikan data, serta manajemen semua role termasuk Admin. Dapat membuat, menghapus, atau mengubah hak akses siapa pun, dan mengelola pengaturan inti aplikasi."},
	{Name: "Admin", Description: "Hak istimewa luas untuk pengelolaan operasional dan pemeliharaan harian. Bertanggung jawab memastikan kelancaran alur kerja. Dapat menambah, mengedit, menghapus Member, serta mengelola proyek, tugas, dan konten. Memantau aktivitas dan memberikan dukungan kepada Member."},
	{Name: "Member", Description: "Pengguna standar yang fokus pada interaksi inti aplikasi. Hak akses terbatas pada tindakan yang diperlukan untuk menjalankan peran dalam proyek atau sistem. Dapat membuat, memperbarui, atau menghapus data mereka sendiri, berpartisipasi dalam proyek yang ditugaskan, dan mengelola profil pribadi."},
}