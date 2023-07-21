package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/funtory-challenge/client/whatsapp"
	"github.com/matinkhosravani/funtory-challenge/domain"
	"github.com/matinkhosravani/funtory-challenge/util"
	"log"
	"net/http"
	"strconv"
)

const QRCodeEventName = "qrcode"

type ConnectWhatsappHandler struct {
	UserRepo domain.UserRepository
}

func (h *ConnectWhatsappHandler) Handle(c *gin.Context) {
	w := c.Writer
	userIDstr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	jid, err := h.UserRepo.GetJIDByUserID(uint(userID))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	client := whatsapp.NewWhatsmeow(h.UserRepo, uint(userID))
	client.NewClient(jid)

	//no need to get qrcode
	if jid != nil {
		client.Connect()
		c.String(http.StatusOK, "already connected")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "SSE not supported")
		return
	}
	util.SetupSSEHeaders(w)

	qrChan := client.GetQRcodeChannel()
	err = client.Connect()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			event, err := util.FormatServerSentEvent(QRCodeEventName, evt.Code)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			_, err = fmt.Fprint(w, event)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
			flusher.Flush()
		} else {
			log.Printf("Login event: %s , user_id : %d \n", evt.Event, userID)
		}
	}
}
