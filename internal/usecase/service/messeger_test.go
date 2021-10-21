package service_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jokius/news-telegram-bot/internal/usecase/service"
	"github.com/jokius/news-telegram-bot/pkg/mocks"
	"github.com/stretchr/testify/require"
)

const (
	testBaseURL = "https://telegram.test.url/"
	userID      = "1"
	token       = "token"
	url         = "https://telegram.test.url/token"
)

func marshalJSON(text string) ([]byte, error) {
	params := struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}{userID, text}

	return json.Marshal(params)
}

func messenger(t *testing.T) (*service.Messenger, *mocks.MockInterfaceClient, *mocks.MockSource) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	client := mocks.NewMockInterfaceClient(mockCtl)
	logger := mocks.NewMockInterfaceLogger(mockCtl)
	source := mocks.NewMockSource(mockCtl)

	newMessenger := service.NewMessenger(token, testBaseURL, client, source, logger)

	return newMessenger, client, source
}

func TestAuth(t *testing.T) {
	t.Parallel()

	serviceMessenger, client, source := messenger(t)

	t.Run("send message to user", func(t *testing.T) {
		t.Parallel()

		authURL := "https://auth_url.example"
		source.EXPECT().AuthURL().Return(authURL).Times(1)
		body, err := marshalJSON("Ссылка для привязки соц сети: " + authURL)
		require.ErrorIs(t, err, nil)
		client.EXPECT().Post(url, body).Times(1)
		serviceMessenger.Auth(userID)
	})
}

func TestURLAdded(t *testing.T) {
	t.Parallel()

	serviceMessenger, client, _ := messenger(t)

	t.Run("send message to user", func(t *testing.T) {
		t.Parallel()

		body, err := marshalJSON("Ссылка на группу добавлена")
		require.ErrorIs(t, err, nil)
		client.EXPECT().Post(url, body).Times(1)
		serviceMessenger.URLAdded(userID)
	})
}

func TestRemovedGroup(t *testing.T) {
	t.Parallel()

	serviceMessenger, client, _ := messenger(t)

	t.Run("send message to user", func(t *testing.T) {
		t.Parallel()

		body, err := marshalJSON("Ссылка на группу удалена")
		require.ErrorIs(t, err, nil)
		client.EXPECT().Post(url, body).Times(1)
		serviceMessenger.RemovedGroup(userID)
	})
}

func TestStartDateUpdated(t *testing.T) {
	t.Parallel()

	serviceMessenger, client, _ := messenger(t)

	t.Run("send message to user", func(t *testing.T) {
		t.Parallel()

		body, err := marshalJSON("Дата начала проверки обновлена")
		require.ErrorIs(t, err, nil)
		client.EXPECT().Post(url, body).Times(1)
		serviceMessenger.StartDateUpdated(userID)
	})
}

func TestGroupList(t *testing.T) {
	t.Parallel()

	serviceMessenger, client, _ := messenger(t)

	t.Run("send message to user", func(t *testing.T) {
		t.Parallel()

		groups := []string{"1", "2"}
		list := strings.Join(groups, "\n")
		body, err := marshalJSON("Список групп:\n" + list)
		require.ErrorIs(t, err, nil)
		client.EXPECT().Post(url, body).Times(1)
		serviceMessenger.GroupList(userID, groups)
	})
}

func TestIncorrectFormat(t *testing.T) {
	t.Parallel()

	serviceMessenger, client, _ := messenger(t)

	t.Run("send message to user add_url", func(t *testing.T) {
		t.Parallel()

		body, err := marshalJSON("Правильный формат: /add_url ссылка на группу")
		require.ErrorIs(t, err, nil)
		client.EXPECT().Post(url, body).Times(1)
		serviceMessenger.IncorrectFormat(userID, "/add_url")
	})

	t.Run("send message to user del_group", func(t *testing.T) {
		t.Parallel()

		body, err := marshalJSON("Правильный формат: /del_group ссылка на группу")
		require.ErrorIs(t, err, nil)
		client.EXPECT().Post(url, body).Times(1)
		serviceMessenger.IncorrectFormat(userID, "/del_group")
	})

	t.Run("send message to user start_date", func(t *testing.T) {
		t.Parallel()

		body, err := marshalJSON("Правильный формат: /start_date dd.mm.yyyy")
		require.ErrorIs(t, err, nil)
		client.EXPECT().Post(url, body).Times(1)
		serviceMessenger.IncorrectFormat(userID, "/start_date")
	})

	t.Run("send message to user unknown", func(t *testing.T) {
		t.Parallel()

		body, err := marshalJSON("Неизвестная команда")
		require.ErrorIs(t, err, nil)
		client.EXPECT().Post(url, body).Times(1)
		serviceMessenger.IncorrectFormat(userID, "unknown")
	})
}

func TestUnknownSource(t *testing.T) {
	t.Parallel()

	serviceMessenger, client, _ := messenger(t)

	t.Run("send message to user", func(t *testing.T) {
		t.Parallel()

		urlStr := "http://unknown.url"
		body, err := marshalJSON("Неизвестный источник: " + urlStr)
		require.ErrorIs(t, err, nil)
		client.EXPECT().Post(url, body).Times(1)
		serviceMessenger.UnknownSource(userID, urlStr)
	})
}

func TestUnknownError(t *testing.T) {
	t.Parallel()

	serviceMessenger, client, _ := messenger(t)

	t.Run("send message to user", func(t *testing.T) {
		t.Parallel()

		errMessage := "some error"
		body, err := marshalJSON("Неизвестная ошибка: " + errMessage)
		require.ErrorIs(t, err, nil)
		client.EXPECT().Post(url, body).Times(1)
		serviceMessenger.UnknownError(userID, errMessage)
	})
}
