package telegram_api

import (
	"net/http"
	"net/url"
)

//go:generate mockery --name HTTPClientPost --output ./mocks
type HTTPClientPost interface {
	PostForm(url string, data url.Values) (*http.Response, error)
}

type WebHookReqBody struct {
	Message Message `json:"message"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Chat struct {
	ID int `json:"id"`
}

type Config struct {
	Port       int    `env:"PORT" envDefault:"3000"`
	Token      string `env:"TOKEN"`
	Caller     int    `env:"CALLER"`
	TimeField  string `env:"TIMEFIELD"`
	TimeFormat string `env:"TIMEFORMAT"`
}
