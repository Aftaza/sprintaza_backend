package model

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EntityProject struct {
	ID       	uint 			`gorm:"primaryKey;column:project_id"` // PK
	Name     	string 			`gorm:"column:name;type:varchar(100);not null"`
	Description string 			`gorm:"column:description;type:text"`

	CreatedAt 	time.Time      `gorm:"column:created_at"`
	UpdatedAt 	time.Time      `gorm:"column:updated_at"`
	DeletedAt 	gorm.DeletedAt `gorm:"index;column:deleted_at"` // Opsi: untuk soft delete
}

// TableName mengembalikan nama nama tabel di database
func (EntityProject) TableName() string {
	return "project"
}

// Di sini kita akan membuat entri EntityPriority dan EntityColumn default untuk project baru ini.
func (ep *EntityProject) AfterCreate(tx *gorm.DB) (err error) {
	// --- Buat Prioritas Default ---
	for _, title := range DefaultPriorityTitles {
		priority := EntityPriority{
			ProjectID: ep.ID,
			Title:     title,
		}
		if createErr := tx.Create(&priority).Error; createErr != nil {
			logrus.Fatalf("Failed to create default priority '%s' for ProjectID %d: %v\n", title, ep.ID, createErr)
			return createErr
		}
	}
	logrus.Infof("Successfully created default priorities for ProjectID %d\n", ep.ID)

	// --- Buat Kolom Default ---
	for i, name := range DefaultColumnNames {
		column := EntityColumn{
			ProjectID: ep.ID,
			Name:      name,
			Order:     i + 1, // Memberikan nilai order yang berurutan (1, 2, 3...)
		}
		if createErr := tx.Create(&column).Error; createErr != nil {
			logrus.Fatalf("Failed to create default column '%s' for ProjectID %d: %v\n", name, ep.ID, createErr)
			return createErr // Mengembalikan error agar GORM membatalkan transaksi
		}
	}
	logrus.Infof("Successfully created default columns for ProjectID %d\n", ep.ID)

	return nil // Mengembalikan nil jika semua berhasil
}