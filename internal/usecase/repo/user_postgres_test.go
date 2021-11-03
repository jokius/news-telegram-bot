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

func buildUserRepo(t *testing.T) (*postgres.Postgres, *repo.UserRepo, dbcleaner.DbCleaner) {
	t.Helper()

	pgURL := os.Getenv("PG_URL_TEST")
	pg, err := postgres.New(pgURL)
	cleaner := dbcleaner.New()
	pgEngine := engine.NewPostgresEngine(pgURL)
	cleaner.SetEngine(pgEngine)

	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}

	userRepo := repo.NewUserRepo(pg)

	return pg, userRepo, cleaner
}

func TestAddGroupByURL(t *testing.T) {
	pg, userRepo, cleaner := buildUserRepo(t)

	t.Run("without user", func(t *testing.T) {
		cleaner.Acquire("users")
		cleaner.Acquire("groups")
		cleaner.Clean("users")
		cleaner.Clean("groups")

		var user entity.User
		pg.Query.Where(&entity.User{TelegramID: userIDUnit}).First(&user)
		assert.Empty(t, user)

		err := userRepo.AddGroupByURL(userID, "vk", "https://example.com/group1")
		assert.ErrorIs(t, err, nil)

		pg.Query.Where(&entity.User{TelegramID: userIDUnit}).First(&user)
		assert.NotEmpty(t, user)

		var group entity.Group
		pg.Query.Where(&entity.Group{UserID: user.ID, SourceName: "vk", Name: "group1"}).First(&group)
		assert.NotEmpty(t, group)

		cleaner.Clean("users")
		cleaner.Clean("groups")
	})

	t.Run("with user", func(t *testing.T) {
		cleaner.Acquire("users")
		cleaner.Acquire("groups")
		cleaner.Clean("users")
		cleaner.Clean("groups")

		timeNow := time.Now()
		user := entity.User{TelegramID: userIDUnit, CreatedAt: timeNow, UpdatedAt: timeNow}
		err := pg.Query.Create(&user).Error
		assert.ErrorIs(t, err, nil)

		err = userRepo.AddGroupByURL(userID, "vk", "https://example.com/group1")
		assert.ErrorIs(t, err, nil)

		var group entity.Group
		pg.Query.Where(&entity.Group{UserID: user.ID, SourceName: "vk", Name: "group1"}).First(&group)
		assert.NotEmpty(t, group)

		cleaner.Clean("users")
		cleaner.Clean("groups")
	})
}

func TestUpdateStartDate(t *testing.T) {
	pg, userRepo, cleaner := buildUserRepo(t)

	t.Run("without user", func(t *testing.T) {
		cleaner.Acquire("users")
		cleaner.Acquire("groups")
		cleaner.Clean("users")
		cleaner.Clean("groups")

		var user entity.User
		pg.Query.Where(&entity.User{TelegramID: userIDUnit}).First(&user)
		assert.Empty(t, user)

		timeNow := time.Now()
		err := userRepo.UpdateStartDate(userID, timeNow)
		assert.ErrorIs(t, err, nil)

		pg.Query.Where(&entity.User{TelegramID: userIDUnit}).First(&user)
		assert.NotEmpty(t, user)

		cleaner.Clean("users")
		cleaner.Clean("groups")
	})

	t.Run("with user", func(t *testing.T) {
		cleaner.Acquire("users")
		cleaner.Acquire("groups")
		cleaner.Clean("users")
		cleaner.Clean("groups")

		timeNow := time.Now()
		user := entity.User{TelegramID: userIDUnit, CreatedAt: timeNow, UpdatedAt: timeNow}
		err := pg.Query.Create(&user).Error
		assert.ErrorIs(t, err, nil)

		group := entity.Group{
			UserID:       user.ID,
			SourceName:   "vk",
			Name:         "group1",
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
			LastUpdateAt: timeNow,
		}
		err = pg.Query.Create(&group).Error
		assert.ErrorIs(t, err, nil)

		dayBefore := timeNow.AddDate(0, 0, -1).UTC()
		err = userRepo.UpdateStartDate(userID, dayBefore)
		assert.ErrorIs(t, err, nil)

		pg.Query.Where(&entity.Group{UserID: user.ID, SourceName: "vk", Name: "group1"}).First(&group)
		assert.NotEmpty(t, group)
		assert.Equal(t, group.LastUpdateAt, dayBefore)

		cleaner.Clean("users")
		cleaner.Clean("groups")
	})
}

func TestRemoveGroup(t *testing.T) {
	pg, userRepo, cleaner := buildUserRepo(t)

	t.Run("without user", func(t *testing.T) {
		cleaner.Acquire("users")
		cleaner.Acquire("groups")
		cleaner.Clean("users")
		cleaner.Clean("groups")

		var user entity.User
		pg.Query.Where(&entity.User{TelegramID: userIDUnit}).First(&user)
		assert.Empty(t, user)

		err := userRepo.RemoveGroup(userID, "vk", "https://example.com/group1")
		assert.ErrorIs(t, err, nil)

		pg.Query.Where(&entity.User{TelegramID: userIDUnit}).First(&user)
		assert.NotEmpty(t, user)

		cleaner.Clean("users")
		cleaner.Clean("groups")
	})

	t.Run("with user", func(t *testing.T) {
		cleaner.Acquire("users")
		cleaner.Acquire("groups")
		cleaner.Clean("users")
		cleaner.Clean("groups")

		timeNow := time.Now()
		user := entity.User{TelegramID: userIDUnit, CreatedAt: timeNow, UpdatedAt: timeNow}
		err := pg.Query.Create(&user).Error
		assert.ErrorIs(t, err, nil)

		group := entity.Group{
			UserID:       user.ID,
			SourceName:   "vk",
			Name:         "group1",
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
			LastUpdateAt: timeNow,
		}
		err = pg.Query.Create(&group).Error
		assert.ErrorIs(t, err, nil)

		err = userRepo.RemoveGroup(userID, "vk", "https://example.com/group1")
		assert.ErrorIs(t, err, nil)

		var emptyGroup entity.Group
		pg.Query.Where(&entity.Group{UserID: user.ID, SourceName: "vk", Name: "group1"}).First(&emptyGroup)
		assert.Empty(t, emptyGroup)

		cleaner.Clean("users")
		cleaner.Clean("groups")
	})
}

func TestGroups(t *testing.T) {
	pg, userRepo, cleaner := buildUserRepo(t)

	t.Run("without user", func(t *testing.T) {
		cleaner.Acquire("users")
		cleaner.Acquire("groups")
		cleaner.Clean("users")
		cleaner.Clean("groups")

		var user entity.User
		pg.Query.Where(&entity.User{TelegramID: userIDUnit}).First(&user)
		assert.Empty(t, user)

		groups, err := userRepo.Groups(userID)
		assert.ErrorIs(t, err, nil)
		assert.Empty(t, groups)

		pg.Query.Where(&entity.User{TelegramID: userIDUnit}).First(&user)
		assert.NotEmpty(t, user)

		cleaner.Clean("users")
		cleaner.Clean("groups")
	})

	t.Run("with user", func(t *testing.T) {
		cleaner.Acquire("users")
		cleaner.Acquire("groups")
		cleaner.Clean("users")
		cleaner.Clean("groups")

		timeNow := time.Now().UTC()
		user := entity.User{TelegramID: userIDUnit, CreatedAt: timeNow, UpdatedAt: timeNow}
		err := pg.Query.Create(&user).Error
		assert.ErrorIs(t, err, nil)

		group := entity.Group{
			UserID:       user.ID,
			SourceName:   "vk",
			Name:         "group1",
			CreatedAt:    timeNow,
			UpdatedAt:    timeNow,
			LastUpdateAt: timeNow,
		}
		err = pg.Query.Create(&group).Error
		assert.ErrorIs(t, err, nil)

		groups, err := userRepo.Groups(userID)
		assert.ErrorIs(t, err, nil)
		assert.NotEmpty(t, groups)
		assert.Equal(t, groups[0], group)

		cleaner.Clean("users")
		cleaner.Clean("groups")
	})
}
