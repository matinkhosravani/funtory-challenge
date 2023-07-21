package domain

type WhatsappClient interface {
	NewClient(jid *string) error
	Connect() error
	GetQRcodeChannel() <-chan QRCodeevent
	AddEventHandler(evtHandler EventHandler)
}
type EventHandler func(evt interface{})

type QRCodeevent struct {
	Code  string
	Event string
}
