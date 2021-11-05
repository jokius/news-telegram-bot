package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jokius/news-telegram-bot/internal/entity"
	"github.com/jokius/news-telegram-bot/internal/usecase/service"
	"github.com/jokius/news-telegram-bot/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func sourceVk(t *testing.T) (*service.VkSource, *mocks.MockInterfaceClient) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	client := mocks.NewMockInterfaceClient(mockCtl)

	newSource := service.NewVkSource("token", client)

	return newSource, client
}

func TestName(t *testing.T) {
	t.Parallel()

	source, _ := sourceVk(t)

	t.Run("get name", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "vk", source.Name())
	})
}

func TestGetGroupMessages(t *testing.T) {
	t.Parallel()

	source, client := sourceVk(t)

	t.Run("get messages", func(t *testing.T) {
		t.Parallel()

		id := "test_id"
		offset := 100
		url := "https://api.vk.com/method/wall.get?v=5.131&count=100&access_token=token&domain=test_id&offset=100"
		client.EXPECT().GetJSON(url, &entity.VkResponse{}).Times(1)
		_, err := source.GetGroupMessages(id, offset)
		assert.ErrorIs(t, err, nil)
	})
}
