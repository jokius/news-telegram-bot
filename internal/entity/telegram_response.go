package entity

type TelegramResult struct {
	Message TelegramMessage `json:"message"`
}

type TelegramMessage struct {
	Text string       `json:"text"`
	User TelegramUser `json:"from"`
}

type TelegramUser struct {
	ID    string `json:"id"`
	IsBot bool   `json:"is_bot"`
}
