package repository

import (
	"github.com/matinkhosravani/funtory-challenge/app"
	"github.com/matinkhosravani/funtory-challenge/domain"
	"github.com/matinkhosravani/funtory-challenge/repository/postgres"
	"log"
)

func NewUserRepository() domain.UserRepository {
	switch app.GetEnv().DBType {
	case "postgres":
		db, err := postgres.NewGorm()
		if err != nil {
			log.Fatal(err.Error())
		}
		return &postgres.UserRepository{DB: db}
	}

	return nil
}
