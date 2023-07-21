package event

import (
	"fmt"
	"github.com/matinkhosravani/funtory-challenge/domain"
	"github.com/matinkhosravani/funtory-challenge/repository"
)

// it could be abstracted but i dont think its a need right now
type whatsappEventHandler struct {
	userRepo domain.UserRepository
}

func NewWhatsappEventHandler() *whatsappEventHandler {
	return &whatsappEventHandler{
		userRepo: repository.NewUserRepository(),
	}
}

func (h *whatsappEventHandler) OnConnect() {
	fmt.Println("connected")
}

func (h *whatsappEventHandler) OnLogOut(userID uint) {
	h.userRepo.SetJIDByUserID(userID, nil)
}

func (h *whatsappEventHandler) OnPairSuccess(userID uint, jid *string) {
	h.userRepo.SetJIDByUserID(userID, jid)
}

func (h *whatsappEventHandler) OnQRScannedWithoutMultidevice() {
	fmt.Println("Multidevice is not enabled")
}
