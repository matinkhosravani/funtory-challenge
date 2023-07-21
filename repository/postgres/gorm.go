package postgres

import (
	"fmt"
	"github.com/matinkhosravani/funtory-challenge/app"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGorm() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		app.GetEnv().DBUser,
		app.GetEnv().DBPass,
		app.GetEnv().DBHost,
		app.GetEnv().DBPort,
		app.GetEnv().DBName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&user{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
