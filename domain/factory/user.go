package factory

import (
	"github.com/bxcodec/faker/v4"
	"github.com/matinkhosravani/funtory-challenge/domain"
	"log"
)

type UserFactory struct {
	userRepo domain.UserRepository
	user     *domain.User
}

// NewDefaultUserFactory Returning Pointer to struct for chaining methods
// it's kinda a form of Builder pattern
func NewDefaultUserFactory(userRepo domain.UserRepository) *UserFactory {
	jid := faker.NAME
	pf := &UserFactory{
		userRepo: userRepo,
		user: &domain.User{
			ID:  1,
			JID: &jid,
		},
	}

	return pf
}
func (f *UserFactory) WithID(id uint) *UserFactory {
	f.user.ID = id
	return f
}

func (f *UserFactory) WithJID(n *string) *UserFactory {
	f.user.JID = n
	return f
}

func (f *UserFactory) Get() *domain.User {
	return f.user
}

func (f *UserFactory) Build() *domain.User {
	err := f.userRepo.Store(f.user)
	if err != nil {
		log.Fatal(err)
	}

	return f.user
}
