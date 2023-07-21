package postgres

import (
	"github.com/matinkhosravani/funtory-challenge/domain"
	"gorm.io/gorm"
)

type user struct {
	ID  uint    `gorm:"primaryKey"`
	JID *string `gorm:"column:jid;nullable"`
}

type UserRepository struct {
	DB *gorm.DB
}

func (r UserRepository) SetJIDByUserID(userID uint, JID *string) error {
	return r.DB.Save(&user{ID: userID, JID: JID}).Error
}

func (r UserRepository) GetJIDByUserID(ID uint) (*string, error) {
	var u user
	err := r.DB.Where("id = ? ", ID).
		First(&u).Error
	if err != nil {
		return nil, err
	}

	return u.JID, nil
}

func (r *UserRepository) Store(domainUser *domain.User) error {
	var p user
	p.JID = domainUser.JID

	err := r.DB.Create(&p).Error
	if err != nil {
		return err
	}
	domainUser.ID = p.ID
	domainUser.JID = p.JID

	return nil
}

func (r *UserRepository) Empty() error {
	return r.DB.Exec("Truncate Table users").Error
}
