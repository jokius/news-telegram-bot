// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"time"

	"github.com/jokius/news-telegram-bot/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=../../pkg/mocks/interfaces_mocks.go -package=mocks

type (
	// User - base user cases.
	User interface {
		TelegramCallback(entity.TelegramResult) error
		AuthToken(id, code string)
	}

	// Messenger - send message to telegram.
	Messenger interface {
		Auth(id string)
		URLAdded(id string)
		RemovedGroup(id string)
		StartDateUpdated(id string)
		GroupList(id string, groups []string)
		IncorrectFormat(id, command string)
		UnknownSource(id, url string)
		UnknownError(id, text string)
	}

	// Source - to work with groups source.
	Source interface {
		Name() string
		AuthURL() (url string)
		GetToken(code string) (token string, err error)
		GetGroupsMessages()
	}

	// UserRepo - user db interaction.
	UserRepo interface {
		AddToken(id, token, sourceName string) (err error)
		AddGroupByURL(id, url string) (err error)
		UpdateStartDate(id string, date time.Time) (err error)
		RemoveGroup(id, url string) (err error)
	}

	// UserWebAPI - user web interaction.
	UserWebAPI interface {
		VkCallback(*entity.User) (entity.User, error)
	}
)
