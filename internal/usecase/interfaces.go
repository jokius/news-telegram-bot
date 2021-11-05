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
		URLAdded(id uint64)
		RemovedGroup(id uint64)
		StartDateUpdated(id uint64)
		GroupList(id uint64, groups []string)
		IncorrectFormat(id uint64, command string)
		UnknownSource(id uint64, url string)
		UnknownError(id uint64, text string)
		Message(id uint64, text string)
	}

	// Source - to work with groups source.
	Source interface {
		Name() string
		GetGroupMessages(id string, offset int) (result entity.VkResult, err error)
	}

	// UserRepo - user db interaction.
	UserRepo interface {
		AddGroupByURL(id uint64, source, url string) (err error)
		UpdateStartDate(id uint64, date time.Time) (err error)
		RemoveGroup(id uint64, source, url string) (err error)
		Groups(id uint64) (groups []entity.Group, err error)
	}

	GroupRepo interface {
		AllBySource(source string) (groups []entity.Group, err error)
		Update(group *entity.Group) (err error)
	}

	MessageRepo interface {
		Add(groupID, messageID uint64, source string, messageAt time.Time) (err error)
		Last(groupID uint64) (message entity.Message)
	}
)
