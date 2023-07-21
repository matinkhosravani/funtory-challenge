package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/funtory-challenge/app"
	"github.com/matinkhosravani/funtory-challenge/domain"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type ConnectWhatsappHandler struct {
	UserRepo domain.UserRepository
	userID   uint
}

func (h *ConnectWhatsappHandler) Handle(c *gin.Context) {
	w := c.Writer
	userIDstr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	h.userID = uint(userID)

	jidText, err := h.UserRepo.GetJIDByUserID(uint(userID))
	if err != nil {
		c.Status(http.StatusBadRequest)
		w.WriteString(err.Error())
		return
	}

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("postgres", "postgresql://root:123456@127.0.0.1:5432/funtory?sslmode=disable", dbLog)
	if err != nil {
		panic(err)
	}
	//no need to get qrcode
	if jidText != nil {
		jid, _ := types.ParseJID(*jidText)
		store, err := container.GetDevice(jid)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		clientLog := waLog.Stdout("Client", "DEBUG", true)
		client := whatsmeow.NewClient(store, clientLog)
		setProxy(client)
		client.AddEventHandler(h.eventHandler)
		client.Connect()
		c.String(http.StatusOK, "already connected")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "SSE not supported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	deviceStore := container.NewDevice()
	if err != nil {
		log.Print(deviceStore)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(h.eventHandler)
	setProxy(client)

	qrChan, _ := client.GetQRChannel(context.Background())
	err = client.Connect()
	if err != nil {
		panic(err)
	}
	for evt := range qrChan {
		if evt.Event == "code" {
			event, err := formatServerSentEvent("qrcode", evt.Code)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			// Render the QR code here
			// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			// or just manually echo 2@... | qrencode -t ansiutf8 in a terminal
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			_, err = fmt.Fprint(w, event)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			flusher.Flush()
		} else {
			log.Println("Login event:", evt.Event)
		}
	}

	fmt.Println("Finished sending price updates...")
}

func setProxy(client *whatsmeow.Client) {
	if app.GetEnv().ProxyEnabled {
		client.SetProxy(func(r *http.Request) (*url.URL, error) {
			u, _ := url.Parse("socks5://127.0.0.1:2080")
			return u, nil
		})
	}
}

func formatServerSentEvent(event string, data any) (string, error) {
	m := map[string]any{
		"data": data,
	}

	buff := bytes.NewBuffer([]byte{})

	encoder := json.NewEncoder(buff)

	err := encoder.Encode(m)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("event: %s\n", event))
	sb.WriteString(fmt.Sprintf("data: %v\n\n", buff.String()))

	return sb.String(), nil
}

func (h ConnectWhatsappHandler) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Connected:
		fmt.Println("connected")
	case *events.LoggedOut:
		h.UserRepo.SetJIDByUserID(h.userID, nil)
	case *events.PairSuccess:
		jid := v.ID.String()
		h.UserRepo.SetJIDByUserID(h.userID, &jid)
	case *events.QRScannedWithoutMultidevice:
		fmt.Println("Multidevice is not enabled")
	}
}
