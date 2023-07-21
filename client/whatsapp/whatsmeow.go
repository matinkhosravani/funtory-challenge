package whatsapp

import (
	"context"
	"fmt"
	"github.com/matinkhosravani/funtory-challenge/app"
	"github.com/matinkhosravani/funtory-challenge/domain"
	wm "go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"net/http"
	"net/url"
)

type Whatsmeow struct {
	client   *wm.Client
	UserRepo domain.UserRepository
	UserID   uint
}

func NewWhatsmeow(userRepo domain.UserRepository, userID uint) domain.WhatsappClient {
	return &Whatsmeow{UserRepo: userRepo, UserID: userID}
}

func (w *Whatsmeow) Connect() error {
	if app.GetEnv().ProxyEnabled {
		w.client.SetProxy(func(r *http.Request) (*url.URL, error) {
			u, _ := url.Parse("socks5://127.0.0.1:2080")
			return u, nil
		})
	}
	return w.client.Connect()
}

func (w *Whatsmeow) AddEventHandler(evtHandler domain.EventHandler) {
	w.client.AddEventHandler(wm.EventHandler(evtHandler))
}

func (w *Whatsmeow) NewClient(jidText *string) error {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("postgres", "postgresql://root:123456@127.0.0.1:5432/funtory?sslmode=disable", dbLog)
	if err != nil {
		return err
	}

	var deviceStore *store.Device
	if jidText == nil {
		deviceStore = container.NewDevice()
	} else {
		jid, _ := types.ParseJID(*jidText)
		deviceStore, err = container.GetDevice(jid)
		if err != nil {
			return err
		}
	}

	clientLog := waLog.Stdout("client", "DEBUG", true)
	w.client = wm.NewClient(deviceStore, clientLog)
	w.client.AddEventHandler(w.eventHandler)

	return err
}

func (w *Whatsmeow) GetQRcodeChannel() <-chan domain.QRCodeevent {
	out := make(chan domain.QRCodeevent)

	qrChan, _ := w.client.GetQRChannel(context.Background())
	w.Connect()
	go func() {
		defer close(out)
		for evt := range qrChan {
			out <- domain.QRCodeevent{
				Code:  evt.Code,
				Event: evt.Event,
			}
		}
	}()

	return out
}

func (w *Whatsmeow) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Connected:
		fmt.Println("connected")
	case *events.LoggedOut:
		w.UserRepo.SetJIDByUserID(w.UserID, nil)
	case *events.PairSuccess:
		jid := v.ID.String()
		w.UserRepo.SetJIDByUserID(w.UserID, &jid)
	case *events.QRScannedWithoutMultidevice:
		fmt.Println("Multidevice is not enabled")
	}
}
