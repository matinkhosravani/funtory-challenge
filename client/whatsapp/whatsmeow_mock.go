package whatsapp

import (
	"github.com/matinkhosravani/funtory-challenge/domain"
)

type WhatsmeowMock struct {
	NewClientFn        func(jid *string) error
	ConnectFn          func() error
	SetUserIDFn        func(userID uint)
	GetQRcodeChannelFn func() <-chan domain.QRCodeevent
	AddEventHandlerFn  func(evtHandler domain.EventHandler)
}

func (w *WhatsmeowMock) NewClient(jid *string) error {
	if w != nil && w.NewClientFn != nil {
		return w.NewClientFn(jid)
	}

	return nil
}

func (w *WhatsmeowMock) SetUserID(userID uint) {
	if w != nil && w.SetUserIDFn != nil {
		w.SetUserIDFn(userID)
	}
}

func (w *WhatsmeowMock) Connect() error {
	if w != nil && w.ConnectFn != nil {
		return w.ConnectFn()
	}

	return nil
}

func (w *WhatsmeowMock) GetQRcodeChannel() <-chan domain.QRCodeevent {
	if w != nil && w.GetQRcodeChannelFn != nil {
		return w.GetQRcodeChannelFn()
	}

	return nil
}

func (w *WhatsmeowMock) AddEventHandler(evtHandler domain.EventHandler) {
	if w != nil && w.AddEventHandlerFn != nil {
		w.AddEventHandlerFn(evtHandler)
	}

	w.AddEventHandlerFn(evtHandler)
}
