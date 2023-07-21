package postgres

import "gorm.io/gorm"

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
