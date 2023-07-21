package domain

type User struct {
	id  uint   `json:"id"`
	jid string `json:"jid"`
}

type UserRepository interface {
	GetJIDByUserID(ID uint) (*string, error)
	SetJIDByUserID(userID uint, JID *string) error
}
