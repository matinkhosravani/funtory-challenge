package domain

type User struct {
	ID  uint   `json:"id"`
	JID string `json:"jid"`
}

type UserRepository interface {
	GetJIDByUserID(ID uint) (*string, error)
	SetJIDByUserID(userID uint, JID *string) error
}
