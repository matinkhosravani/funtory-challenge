package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/funtory-challenge/domain"
	"github.com/matinkhosravani/funtory-challenge/util"
	"log"
	"net/http"
	"strconv"
)

const QRCodeEventName = "qrcode"

type ConnectWhatsappHandler struct {
	UserRepo domain.UserRepository
	Client   domain.WhatsappClient
}

// Handle is the handler function for connecting user's whatsapp account
// @Summary connect whatsapp account to database
// @Description is the handler function for connecting user's whatsapp account.
// if user has already connected his account then there is no need to generate Qrcode
// but if there is no jid assigned him in database then server  opens a connection
// and push qrcode to the client once the previous qrcode expires
// @ID connect
// @Param id path int true "id of user"
// @Success 200 {string} example
// @Router /connect/{id} [get]
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
	h.Client.SetUserID(uint(userID))
	h.Client.NewClient(jid)

	//no need to get qrcode
	if jid != nil {
		h.Client.Connect()
		c.String(http.StatusOK, "already connected")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		c.String(http.StatusInternalServerError, "SSE not supported")
		return
	}
	util.SetupSSEHeaders(w)
	qrChan := h.Client.GetQRcodeChannel()
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
