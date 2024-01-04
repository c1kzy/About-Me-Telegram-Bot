package telegram_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"telegram/internal"

	"github.com/phuslu/log"
)

var (
	lock      = &sync.Mutex{}
	singleApi *API
)

type API struct {
	Client HTTPClientPost
	Url    string
}

func GetAPI(cfg *Config) *API {
	if singleApi == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleApi == nil {
			singleApi = &API{
				Client: http.DefaultClient,
				Url:    fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", cfg.Token),
			}
			log.Info().Msg("API created")
		}
	}

	return singleApi
}

// Message to send
func (api *API) sendResponse(chatID int, text string) error {
	response, err := api.Client.PostForm(api.Url, url.Values{"chat_id": {strconv.Itoa(chatID)}, "text": {text}})
	if err != nil {
		return fmt.Errorf("sending response failed. ChatID:%v, Text:%v.Error:%v", chatID, text, err)
	}

	if response == nil {
		return fmt.Errorf("response for ChatID:%v is nil. Input text:%v", chatID, text)
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 && response.StatusCode < 500 {
		responseBody, readErr := io.ReadAll(response.Body)
		log.Warn().Msgf("Unable to send response for ChatID:%v. Text:%v. Response body:%v. Response error: %v", chatID, text, string(responseBody), readErr)

	}

	if response.StatusCode >= 500 {
		return fmt.Errorf("%v internal server error. ChatID:%v, Text:%s", response.StatusCode, chatID, text)
	}

	log.Debug().Msgf("Response to ChatID:%v sent. Message:%v. Response: %v", chatID, text, response)

	return nil
}

// Handler function
func (api *API) TelegramHandler(_ http.ResponseWriter, r *http.Request) {
	var body *WebHookReqBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error().Err(err).Msgf("error occurred decoding message body: %v. Error:%v", r.Body, err)
	}

	// Different cases handling
	input := body.Message.Text

	textToSend := "Unrecognized command. Type /help for more information"

	switch input {
	case "/start":
		textToSend = internal.StartText

	case "/about":
		textToSend = internal.AboutMeText

	case "/links":
		textToSend = internal.LinkText

	case "/help":
		textToSend = internal.HelpText

	}

	err := api.sendResponse(body.Message.Chat.ID, textToSend)
	if err != nil {
		log.Error().Err(fmt.Errorf("sendResponse error: %w", err))
		api.sendResponse(body.Message.Chat.ID, fmt.Sprintf("error sending reply for /start command %v", err))
		return
	}
}
