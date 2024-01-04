package telegram_api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/phuslu/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (c *MockHTTPClient) PostForm(url string, data url.Values) (*http.Response, error) {
	args := c.Called(url, data)
	return args.Get(0).(*http.Response), args.Error(1)
}

func marshJSON(t *testing.T, text string, chatID int) string {
	var reqBody *WebHookReqBody

	reqBody = &WebHookReqBody{
		Message: Message{
			Text: text,
			Chat: Chat{
				ID: chatID,
			},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func TestHandlerTelegram(t *testing.T) {
	var cfg Config

	log.DefaultLogger = log.Logger{
		Level:      log.DebugLevel,
		Caller:     cfg.Caller,
		TimeField:  cfg.TimeField,
		TimeFormat: time.RFC850,
		Writer:     &log.ConsoleWriter{},
	}

	path, err := filepath.Abs("../.env")
	if err != nil {
		log.Fatal().Err(err).Msgf("error getting .env file path.Details:%v", err)
	}
	log.Debug().Msgf("Path to .env: %v", path)

	if envErr := godotenv.Load(os.ExpandEnv(path)); envErr != nil {
		log.Fatal().Msgf("Error loading .env file %v", envErr)
	}

	if parseErr := env.Parse(&cfg); err != nil {
		log.Error().Msgf("parsing env failed. Details:%v", parseErr)
	}

	recorder := httptest.NewRecorder()

	api := GetAPI(&cfg)
	log.Debug().Msgf("Config loaded: %v", cfg)

	tests := []struct {
		name    string
		request *http.Request
	}{
		{
			name:    "/start",
			request: httptest.NewRequest(http.MethodPost, "/telegram", bytes.NewBuffer([]byte(marshJSON(t, "/start", 358383178)))),
		},
		{
			name:    "/about",
			request: httptest.NewRequest(http.MethodPost, "/telegram", bytes.NewBuffer([]byte(marshJSON(t, "/about", 358383178)))),
		},
		{
			name:    "/links",
			request: httptest.NewRequest(http.MethodPost, "/telegram", bytes.NewBuffer([]byte(marshJSON(t, "/links", 358383178)))),
		},
		{
			name:    "/help",
			request: httptest.NewRequest(http.MethodPost, "/telegram", bytes.NewBuffer([]byte(marshJSON(t, "/help", 358383178)))),
		},

		{
			name:    "invalid command",
			request: httptest.NewRequest(http.MethodPost, "/telegram", bytes.NewBuffer([]byte(marshJSON(t, "/invalid", 358383178)))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api.TelegramHandler(recorder, tt.request)
			assert.Equal(t, http.StatusOK, recorder.Code)
		})
	}
}
