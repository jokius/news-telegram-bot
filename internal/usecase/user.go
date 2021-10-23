// Package usecase implement User interface
package usecase

import (
	"fmt"
	"strings"
	"time"

	"github.com/jokius/news-telegram-bot/internal/entity"
	"github.com/jokius/news-telegram-bot/pkg/errors"
)

// UserUseCase -.
type UserUseCase struct {
	repo   UserRepo
	msg    Messenger
	source Source
}

const (
	commandWithParams = 2
)

// NewUserUseCase - init.
func NewUserUseCase(r UserRepo, m Messenger, s Source) *UserUseCase {
	return &UserUseCase{r, m, s}
}

// TelegramCallback - parse telegram callback.
func (uc *UserUseCase) TelegramCallback(telegramResult entity.TelegramResult) (err error) {
	message := telegramResult.Message
	user := message.User

	if user.IsBot {
		return fmt.Errorf("%w", errors.ErrBotMessage)
	}

	id := user.ID
	text := message.Text
	textSlice := strings.Split(text, " ")

	switch {
	case text == "/list":
		uc.groupList(id)
	case len(textSlice) >= commandWithParams:
		uc.messageWithParams(textSlice, id)
	default:
		uc.msg.IncorrectFormat(id, textSlice[0])
	}

	return
}

func (uc *UserUseCase) groupList(id string) {
	groups, err := uc.repo.Groups(id)
	if err != nil {
		uc.msg.UnknownError(id, "something wrong: "+err.Error())

		return
	}

	list := make([]string, len(groups))
	for i := range groups {
		list[i] = groups[i].GroupLink
	}

	uc.msg.GroupList(id, list)
}

func (uc *UserUseCase) messageWithParams(textSlice []string, id string) {
	text := textSlice[1]

	switch textSlice[0] {
	case "/add_url":
		uc.addURL(id, text)
	case "/start_date":
		uc.startDate(id, text)
	case "/del_group":
		uc.removeGroup(id, text)
	default:
		uc.msg.IncorrectFormat(id, "unknown")
	}
}

func (uc *UserUseCase) addURL(id, text string) {
	err := uc.repo.AddGroupByURL(id, text)
	if err == nil {
		uc.msg.URLAdded(id)
	} else {
		uc.errBD(id, err)
	}
}

func (uc *UserUseCase) removeGroup(id, text string) {
	err := uc.repo.RemoveGroup(id, text)
	if err == nil {
		uc.msg.RemovedGroup(id)
	} else {
		uc.errBD(id, err)
	}
}

func (uc *UserUseCase) startDate(id, text string) {
	t, err := time.Parse("02.01.2006", text)
	if err != nil {
		uc.msg.IncorrectFormat(id, "start_date")

		return
	}

	err = uc.repo.UpdateStartDate(id, t)
	if err == nil {
		uc.msg.StartDateUpdated(id)
	} else {
		uc.errBD(id, err)
	}
}

func (uc *UserUseCase) errBD(id string, err error) {
	uc.msg.UnknownError(id, "something wrong: "+err.Error())
}
