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

func buildGroupRepo(t *testing.T) (*postgres.Postgres, *repo.GroupRepo, dbcleaner.DbCleaner) {
	t.Helper()

	pgURL := os.Getenv("PG_URL_TEST")
	pg, err := postgres.New(pgURL)
	cleaner := dbcleaner.New()
	pgEngine := engine.NewPostgresEngine(pgURL)
	cleaner.SetEngine(pgEngine)

	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}

	groupRepo := repo.NewGroupRepo(pg)

	return pg, groupRepo, cleaner
}

func TestAllBySource(t *testing.T) {
	pg, groupRepo, cleaner := buildGroupRepo(t)

	t.Run("run", func(t *testing.T) {
		cleaner.Acquire("users")
		cleaner.Acquire("groups")
		cleaner.Clean("users")
		cleaner.Clean("groups")

		timeNow := time.Now().UTC()
		user := entity.User{TelegramID: userID, CreatedAt: timeNow, UpdatedAt: timeNow}
		err := pg.Query.Create(&user).Error
		assert.ErrorIs(t, err, nil)

		group := entity.Group{UserID: user.ID, SourceName: "vk", Name: "test_group", LastUpdateAt: timeNow, CreatedAt: timeNow, UpdatedAt: timeNow}
		otherGroup := entity.Group{UserID: user.ID, SourceName: "other", Name: "test_group", LastUpdateAt: timeNow, CreatedAt: timeNow, UpdatedAt: timeNow}
		err = pg.Query.Create(&group).Error
		assert.ErrorIs(t, err, nil)

		err = pg.Query.Create(&otherGroup).Error
		assert.ErrorIs(t, err, nil)

		groups, err := groupRepo.AllBySource("vk")
		assert.ErrorIs(t, err, nil)
		assert.NotEmpty(t, groups)
		assert.Equal(t, len(groups), 1)
		assert.NotEmpty(t, groups[0].User)

		cleaner.Clean("users")
		cleaner.Clean("groups")
	})
}

func TestUpdateGroup(t *testing.T) {
	pg, groupRepo, cleaner := buildGroupRepo(t)

	t.Run("run", func(t *testing.T) {
		cleaner.Acquire("users")
		cleaner.Acquire("groups")
		cleaner.Clean("users")
		cleaner.Clean("groups")

		timeNow := time.Now().UTC()
		user := entity.User{TelegramID: userID, CreatedAt: timeNow, UpdatedAt: timeNow}
		err := pg.Query.Create(&user).Error
		assert.ErrorIs(t, err, nil)

		group := entity.Group{UserID: user.ID, SourceName: "vk", Name: "test_group", LastUpdateAt: timeNow, CreatedAt: timeNow, UpdatedAt: timeNow}
		err = pg.Query.Create(&group).Error
		assert.ErrorIs(t, err, nil)

		dayBefore := timeNow.AddDate(0, 0, -1).UTC()
		group.LastUpdateAt = dayBefore
		err = groupRepo.Update(&group)
		assert.ErrorIs(t, err, nil)

		var updatedGroup entity.Group
		err = pg.Query.Where(&entity.Group{ID: group.ID}).First(&updatedGroup).Error
		assert.ErrorIs(t, err, nil)
		assert.Equal(t, updatedGroup.LastUpdateAt, group.LastUpdateAt)

		cleaner.Clean("users")
		cleaner.Clean("groups")
	})
}
