package repo

import (
	"time"

	"github.com/jokius/news-telegram-bot/internal/entity"
	"github.com/jokius/news-telegram-bot/pkg/postgres"
)

type MessageRepo struct {
	db *postgres.Postgres
}

func NewMessageRepo(pg *postgres.Postgres) *MessageRepo {
	return &MessageRepo{pg}
}

func (m MessageRepo) Add(groupID uint64, source, messageID string, messageAt time.Time) error {
	t := time.Now()
	message := entity.Message{
		GroupID:   groupID,
		MessageID: messageID,
		Source:    source,
		MessageAt: messageAt,
		CreatedAt: t,
		UpdatedAt: t,
	}

	return m.db.Query.Create(&message).Error
}

func (m MessageRepo) Last(groupID uint64) (message entity.Message) {
	m.db.Query.Where(&entity.Message{GroupID: groupID}).Order("message_at desc").First(&message)

	return
}
