package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/matinkhosravani/funtory-challenge/app"
	"github.com/matinkhosravani/funtory-challenge/client/whatsapp"
	"github.com/matinkhosravani/funtory-challenge/domain"
	"github.com/matinkhosravani/funtory-challenge/domain/factory"
	"github.com/matinkhosravani/funtory-challenge/repository"
	"github.com/matinkhosravani/funtory-challenge/repository/mock"
	"github.com/matinkhosravani/funtory-challenge/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestUnit_Handle(t *testing.T) {

	t.Run("user already connected", func(t *testing.T) {
		userID := 123
		jid := "dummy_jid"
		userRepo := &mock.UserRepository{
			GetJIDByUserIDFn: func(ID uint) (*string, error) {
				return &jid, nil
			},
		}
		whatsappClient := &whatsapp.WhatsmeowMock{}

		router := gin.Default()
		handler := &ConnectWhatsappHandler{
			UserRepo: userRepo,
			Client:   whatsappClient,
		}
		router.GET("/connect/:user_id", handler.Handle)
		req, _ := http.NewRequest(http.MethodGet, "/connect/"+strconv.Itoa(userID), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "already connected", resp.Body.String())
	})

	t.Run("user Get qrcode when he is not connected before", func(t *testing.T) {
		userID := 123
		dummyQrCodeevent := domain.QRCodeevent{
			Code:  "dummy qrcode",
			Event: "code",
		}
		userRepo := &mock.UserRepository{}
		whatsappClient := &whatsapp.WhatsmeowMock{
			GetQRcodeChannelFn: func() <-chan domain.QRCodeevent {
				ch := make(chan domain.QRCodeevent, 1)
				ch <- dummyQrCodeevent
				close(ch)
				return ch
			},
		}

		handler := &ConnectWhatsappHandler{
			UserRepo: userRepo,
			Client:   whatsappClient,
		}
		router := gin.Default()
		router.GET("/connect/:user_id", handler.Handle)
		req, _ := http.NewRequest(http.MethodGet, "/connect/"+strconv.Itoa(userID), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
		// Check the response headers for sse
		expectedHeaders := map[string]string{
			"Content-Type":  "text/event-stream",
			"Cache-Control": "no-cache",
			"Connection":    "keep-alive",
		}
		for key, value := range expectedHeaders {
			if resp.Header().Get(key) != value {
				t.Errorf("Expected header %s: %s, got %s", key, value, resp.Header().Get(key))
			}
		}
		data, _ := util.FormatServerSentEvent(QRCodeEventName, dummyQrCodeevent.Code)
		assert.Equal(t, data, resp.Body.String())
	})
}

func TestIntegrate_Handle(t *testing.T) {

	t.Run("user already connected", func(t *testing.T) {
		app.BootTestApp()
		userRepo := repository.NewUserRepository()
		_ = userRepo.Empty()
		jid := "dummy"
		user := factory.NewDefaultUserFactory(userRepo).
			WithJID(&jid).
			Build()
		whatsappClient := &whatsapp.WhatsmeowMock{}
		router := gin.Default()
		handler := &ConnectWhatsappHandler{
			UserRepo: userRepo,
			Client:   whatsappClient,
		}
		router.GET("/connect/:user_id", handler.Handle)
		req, _ := http.NewRequest(http.MethodGet, "/connect/"+strconv.Itoa(int(user.ID)), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "already connected", resp.Body.String())
	})
}
