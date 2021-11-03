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
	}

	// Messenger - send message to telegram.
	Messenger interface {
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
		GetGroupMessages()
	}

	// UserRepo - user db interaction.
	UserRepo interface {
		AddGroupByURL(id, source, url string) (err error)
		UpdateStartDate(id string, date time.Time) (err error)
		RemoveGroup(id, source, url string) (err error)
		Groups(id string) (groups []entity.Group, err error)
	}

	GroupRepo interface {
		AllBySource(source string) (groups []entity.Group, err error)
		Update(group *entity.Group) (err error)
	}
)
