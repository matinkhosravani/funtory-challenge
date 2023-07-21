package mock

type UserRepository struct {
	GetJIDByUserIDFn func(ID uint) (*string, error)
	SetJIDByUserIDFn func(userID uint, JID *string) error
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
