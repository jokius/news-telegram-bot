package repo_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jokius/news-telegram-bot/internal/entity"
	"github.com/jokius/news-telegram-bot/internal/usecase/repo"
	"github.com/jokius/news-telegram-bot/pkg/postgres"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
)

func buildMessageRepo(t *testing.T) (*postgres.Postgres, *repo.MessageRepo, dbcleaner.DbCleaner) {
	t.Helper()

	pgURL := os.Getenv("PG_URL_TEST")
	pg, err := postgres.New(pgURL)
	cleaner := dbcleaner.New()
	pgEngine := engine.NewPostgresEngine(pgURL)
	cleaner.SetEngine(pgEngine)

	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}

	groupRepo := repo.NewMessageRepo(pg)

	return pg, groupRepo, cleaner
}

func TestAddMessage(t *testing.T) {
	pg, messageRepo, cleaner := buildMessageRepo(t)

	t.Run("run", func(t *testing.T) {
		cleaner.Acquire("messages")
		cleaner.Clean("messages")

		messageAt := time.Now().UTC()
		err := messageRepo.Add(groupID, messageID, "vk", messageAt)
		assert.ErrorIs(t, err, nil)

		var message entity.Message
		err = pg.Query.
			Where(&entity.Message{GroupID: groupID, MessageID: messageID, Source: "vk", MessageAt: messageAt}).
			First(&message).
			Error
		assert.ErrorIs(t, err, nil)
		assert.NotEmpty(t, message)

		cleaner.Clean("messages")
	})
}

func TestLastMessage(t *testing.T) {
	pg, messageRepo, cleaner := buildMessageRepo(t)

	t.Run("run", func(t *testing.T) {
		cleaner.Acquire("messages")
		cleaner.Clean("messages")

		messageAt := time.Now().UTC()

		dayBefore := time.Now().AddDate(0, 0, -1).UTC()
		err := pg.Query.Create(&entity.Message{
			GroupID:   groupID,
			MessageID: 2,
			Source:    "vk",
			MessageAt: dayBefore,
			CreatedAt: messageAt,
			UpdatedAt: messageAt,
		}).Error
		assert.ErrorIs(t, err, nil)

		message := entity.Message{
			GroupID:   groupID,
			MessageID: messageID,
			Source:    "vk",
			MessageAt: messageAt,
			CreatedAt: messageAt,
			UpdatedAt: messageAt,
		}

		err = pg.Query.Create(&message).Error
		assert.ErrorIs(t, err, nil)

		lastMessage := messageRepo.Last(groupID)
		assert.Equal(t, message, lastMessage)

		cleaner.Clean("messages")
	})
}
