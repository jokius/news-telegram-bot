package service

import (
	"strconv"

	"github.com/jokius/news-telegram-bot/internal/entity"
	"github.com/jokius/news-telegram-bot/pkg/httpclient"
)

type VkSource struct {
	name   string
	token  string
	client httpclient.InterfaceClient
}

const (
	baseURL = "https://api.vk.com/method/wall.get?v=5.131&count=100"
)

func NewVkSource(token string, client httpclient.InterfaceClient) *VkSource {
	return &VkSource{"vk", token, client}
}

func (v *VkSource) Name() string {
	return v.name
}

func (v *VkSource) GetGroupMessages(id string, offset int) (result entity.VkResult, err error) {
	url := baseURL +
		"&access_token=" + v.token +
		"&posts=" + id +
		"&offset=" + strconv.Itoa(offset)

	err = v.client.GetJSON(url, result)

	return
}
