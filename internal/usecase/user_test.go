package usecase_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jokius/news-telegram-bot/internal/entity"
	"github.com/jokius/news-telegram-bot/internal/usecase"
	"github.com/jokius/news-telegram-bot/pkg/errors"
	"github.com/jokius/news-telegram-bot/pkg/mocks"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const (
	userID   = "1"
	timeText = "10.11.2021"
)

func telegramResult(text string) entity.TelegramResult {
	return entity.TelegramResult{
		Message: entity.TelegramMessage{
			Text: text,
			User: entity.TelegramUser{
				ID:    userID,
				IsBot: false,
			},
		},
	}
}

func user(t *testing.T) (*usecase.UserUseCase, *mocks.MockMessenger, *mocks.MockUserRepo) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	repo := mocks.NewMockUserRepo(mockCtl)
	messenger := mocks.NewMockMessenger(mockCtl)
	source := mocks.NewMockSource(mockCtl)

	newUser := usecase.NewUserUseCase(repo, messenger, source)

	return newUser, messenger, repo
}

func TestTelegramCallback_correct(t *testing.T) {
	t.Parallel()

	userCase, message, repo := user(t)

	t.Run("when add_url", func(t *testing.T) {
		t.Parallel()

		repo.EXPECT().AddGroupByURL(userID, "https://example.com/1").Return(nil).Times(1)
		message.EXPECT().URLAdded(userID).Return().Times(1)
		err := userCase.TelegramCallback(telegramResult("/add_url https://example.com/1"))
		require.ErrorIs(t, err, nil)
	})

	t.Run("when start_date", func(t *testing.T) {
		t.Parallel()

		timeParse, err := time.Parse("02.01.2006", timeText)
		require.ErrorIs(t, err, nil)

		repo.EXPECT().UpdateStartDate(userID, timeParse).Return(nil).Times(1)
		message.EXPECT().StartDateUpdated(userID).Return().Times(1)
		err = userCase.TelegramCallback(telegramResult("/start_date " + timeText))
		require.ErrorIs(t, err, nil)
	})

	t.Run("when del_group", func(t *testing.T) {
		t.Parallel()

		repo.EXPECT().RemoveGroup(userID, "https://example.com/1").Return(nil).Times(1)
		message.EXPECT().RemovedGroup(userID).Return().Times(1)
		err := userCase.TelegramCallback(telegramResult("/del_group https://example.com/1"))
		require.ErrorIs(t, err, nil)
	})

	t.Run("when list", func(t *testing.T) {
		t.Parallel()

		listStr := []string{"1"}
		listGroups := []entity.Group{{GroupLink: "1"}}
		repo.EXPECT().Groups(userID).Return(listGroups, nil).Times(1)
		message.EXPECT().GroupList(userID, listStr).Times(1)
		err := userCase.TelegramCallback(telegramResult("/list"))
		require.ErrorIs(t, err, nil)
	})
}

func TestTelegramCallback_with_db_error(t *testing.T) {
	t.Parallel()

	errBD := gorm.ErrInvalidValue
	userCase, message, repo := user(t)

	t.Run("when add_url", func(t *testing.T) {
		t.Parallel()

		repo.EXPECT().AddGroupByURL(userID, "https://example.com/1").Return(errBD).Times(1) // any error
		message.EXPECT().UnknownError(userID, "something wrong: "+errBD.Error()).Return().Times(1)
		err := userCase.TelegramCallback(telegramResult("/add_url https://example.com/1"))
		require.ErrorIs(t, err, nil)
	})

	t.Run("when start_date", func(t *testing.T) {
		t.Parallel()

		timeParse, err := time.Parse("02.01.2006", timeText)
		require.ErrorIs(t, err, nil)

		repo.EXPECT().UpdateStartDate(userID, timeParse).Return(errBD).Times(1)
		message.EXPECT().UnknownError(userID, "something wrong: "+errBD.Error()).Return().Times(1)
		err = userCase.TelegramCallback(telegramResult("/start_date " + timeText))
		require.ErrorIs(t, err, nil)
	})

	t.Run("when del_group", func(t *testing.T) {
		t.Parallel()

		repo.EXPECT().RemoveGroup(userID, "https://example.com/1").Return(errBD).Times(1) // any error
		message.EXPECT().UnknownError(userID, "something wrong: "+errBD.Error()).Return().Times(1)
		err := userCase.TelegramCallback(telegramResult("/del_group https://example.com/1"))
		require.ErrorIs(t, err, nil)
	})

	t.Run("when list", func(t *testing.T) {
		t.Parallel()

		repo.EXPECT().Groups(userID).Return([]entity.Group{}, errBD).Times(1) // any error
		message.EXPECT().UnknownError(userID, "something wrong: "+errBD.Error()).Return().Times(1)
		err := userCase.TelegramCallback(telegramResult("/list"))
		require.ErrorIs(t, err, nil)
	})
}

func TestTelegramCallback_with_error_noParams(t *testing.T) {
	t.Parallel()

	userCase, message, _ := user(t)

	t.Run("when add_url", func(t *testing.T) {
		t.Parallel()

		message.EXPECT().IncorrectFormat(userID, "/add_url").Return().Times(1)
		err := userCase.TelegramCallback(telegramResult("/add_url"))
		require.ErrorIs(t, err, nil)
	})

	t.Run("when start_date", func(t *testing.T) {
		t.Parallel()

		message.EXPECT().IncorrectFormat(userID, "/start_date").Return().Times(1)
		err := userCase.TelegramCallback(telegramResult("/start_date"))
		require.ErrorIs(t, err, nil)
	})

	t.Run("when del_group", func(t *testing.T) {
		t.Parallel()

		message.EXPECT().IncorrectFormat(userID, "/del_group").Return().Times(1)
		err := userCase.TelegramCallback(telegramResult("/del_group"))
		require.ErrorIs(t, err, nil)
	})
}

func TestTelegramCallback_with_error_other(t *testing.T) {
	t.Parallel()

	userCase, message, _ := user(t)

	t.Run("when is bot", func(t *testing.T) {
		t.Parallel()

		res := telegramResult("")
		res.Message.User.IsBot = true
		err := userCase.TelegramCallback(res)
		require.ErrorIs(t, err, errors.ErrBotMessage)
	})

	t.Run("when start_date", func(t *testing.T) {
		t.Parallel()

		message.EXPECT().IncorrectFormat(userID, "start_date").Return().Times(1)
		err := userCase.TelegramCallback(telegramResult("/start_date no_time"))
		require.ErrorIs(t, err, nil)
	})

	t.Run("when incorrect_command params", func(t *testing.T) {
		t.Parallel()

		message.EXPECT().IncorrectFormat(userID, "unknown").Return().Times(1)
		err := userCase.TelegramCallback(telegramResult("/incorrect_command 123"))
		require.ErrorIs(t, err, nil)
	})
}
