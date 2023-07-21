package mock

import "github.com/matinkhosravani/funtory-challenge/domain"

type UserRepository struct {
	GetJIDByUserIDFn func(ID uint) (*string, error)
	SetJIDByUserIDFn func(userID uint, JID *string) error
	StoreFn          func(user *domain.User) error
}

func (u *UserRepository) GetJIDByUserID(ID uint) (*string, error) {
	if u != nil && u.GetJIDByUserIDFn != nil {
		return u.GetJIDByUserIDFn(ID)
	}

	return nil, nil
}

func (u *UserRepository) SetJIDByUserID(userID uint, JID *string) error {
	if u != nil && u.SetJIDByUserIDFn != nil {
		return u.SetJIDByUserID(userID, JID)
	}

	return nil
}

func (u *UserRepository) Store(user *domain.User) error {
	if u != nil && u.StoreFn != nil {
		return u.StoreFn(user)
	}

	return nil
}

func (u *UserRepository) Empty() error {
	return nil
}
