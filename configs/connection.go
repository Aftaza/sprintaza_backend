package config

import (
	"os"

	model "github.com/Aftaza/sprintaza_backend/models"
	util "github.com/Aftaza/sprintaza_backend/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connection() *gorm.DB {
	databaseURI := make(chan string, 1)

	if os.Getenv("GO_ENV") != "production" {
		databaseURI <- util.GodotEnv("DATABASE_URI_DEV")
	} else {
		databaseURI <- os.Getenv("DATABASE_URI_PROD")
	}

	db, err := gorm.Open(postgres.Open(<-databaseURI), &gorm.Config{})

	if err != nil {
		defer logrus.Info("Connection to Database Failed")
		logrus.Fatal(err.Error())
	}

	if os.Getenv("GO_ENV") != "production" {
		logrus.Info("Connection to Database Successfully")
	}

	err = db.AutoMigrate(
		&model.EntityAchievement{},
		&model.EntityRole{},
		&model.EntityColumn{},
		&model.EntityPriority{},
		&model.EntityUsers{},
		&model.EntityUserXP{},
		&model.EntityUserAchievement{},
		&model.EntityProject{},
		&model.EntityProjectMember{},
		&model.EntityTask{},
		&model.EntitySubtask{},
	)

	if err != nil {
		logrus.Fatal(err.Error())
	}

	SeedAchievements(db)

	return db
}