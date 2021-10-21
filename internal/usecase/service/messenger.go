// Package service implement Messenger interface
package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jokius/news-telegram-bot/internal/usecase"
	"github.com/jokius/news-telegram-bot/pkg/httpclient"
	"github.com/jokius/news-telegram-bot/pkg/logger"
)

// Messenger - messenger to telegram.
type Messenger struct {
	baseURL string
	token   string
	client  httpclient.InterfaceClient
	source  usecase.Source
	logger  logger.InterfaceLogger
}

const (
	_defaultBaseURL = "https://api.telegram.org/"
)

func NewMessenger(token, baseURL string,
	client httpclient.InterfaceClient,
	source usecase.Source,
	l logger.InterfaceLogger) *Messenger {
	if baseURL == "" {
		baseURL = _defaultBaseURL
	} else if lastCh := baseURL[len(baseURL)-1:]; lastCh != "/" {
		baseURL += "/"
	}

	return &Messenger{baseURL, token, client, source, l}
}

func (m *Messenger) Auth(id string) {
	m.sendMessage(id, "Ссылка для привязки соц сети: "+m.source.AuthURL())
}

func (m *Messenger) URLAdded(id string) {
	m.sendMessage(id, "Ссылка на группу добавлена")
}

func (m *Messenger) RemovedGroup(id string) {
	m.sendMessage(id, "Ссылка на группу удалена")
}

func (m *Messenger) StartDateUpdated(id string) {
	m.sendMessage(id, "Дата начала проверки обновлена")
}

func (m *Messenger) GroupList(id string, groups []string) {
	list := strings.Join(groups, "\n")
	m.sendMessage(id, "Список групп:\n"+list)
}

func (m *Messenger) IncorrectFormat(id, command string) {
	var text string

	switch command {
	case "/add_url":
		text = "Правильный формат: /add_url ссылка на группу"
	case "/del_group":
		text = "Правильный формат: /del_group ссылка на группу"
	case "/start_date":
		text = "Правильный формат: /start_date dd.mm.yyyy"
	default:
		text = "Неизвестная команда"
	}

	m.sendMessage(id, text)
}

func (m *Messenger) UnknownSource(id, url string) {
	m.sendMessage(id, "Неизвестный источник: "+url)
}

func (m *Messenger) UnknownError(id, text string) {
	m.sendMessage(id, "Неизвестная ошибка: "+text)
}

func (m *Messenger) sendMessage(id, message string) {
	url := m.baseURL + m.token
	params := struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}{id, message}

	body, err := json.Marshal(params)
	if err != nil {
		m.logger.Fatal(fmt.Errorf("something wrong: %w", err))
	}

	if _, err = m.client.Post(url, body); err != nil {
		m.logger.Fatal(fmt.Errorf("something wrong: %w", err))
	}
}
