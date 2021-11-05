package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jokius/news-telegram-bot/internal/entity"
	"github.com/jokius/news-telegram-bot/internal/usecase"
	"github.com/jokius/news-telegram-bot/pkg/logger"
)

type GrabberVk struct {
	sleep       time.Duration
	source      usecase.Source
	messenger   usecase.Messenger
	groupRepo   usecase.GroupRepo
	messageRepo usecase.MessageRepo
	l           logger.InterfaceLogger
}

func NewVkGrabber(sleep time.Duration, source usecase.Source, messenger usecase.Messenger, groupRepo usecase.GroupRepo,
	messageRepo usecase.MessageRepo, l logger.InterfaceLogger) GrabberVk {
	return GrabberVk{
		sleep:       sleep,
		source:      source,
		messenger:   messenger,
		groupRepo:   groupRepo,
		messageRepo: messageRepo,
		l:           l,
	}
}

func (g *GrabberVk) Start(shutdown chan bool) {
	go func() {
		for {
			select {
			case <-shutdown:
				break
			default:
			}

			go g.grab()
			time.Sleep(g.sleep)
		}
	}()
}

func (g *GrabberVk) grab() {
	groups, err := g.groupRepo.AllBySource(g.source.Name())
	if err != nil {
		g.l.Error(fmt.Errorf("`g.grab` something wrong: %w", err))

		return
	}

	t := time.Now().UTC()

	for i := range groups {
		group := &groups[i]
		if !group.LastUpdateAt.After(t) {
			if err = g.grabGroup(group, t); err != nil {
				return
			}
		}
	}
}

func (g *GrabberVk) grabGroup(group *entity.Group, t time.Time) (err error) {
	err = g.newMessages(group)
	if err != nil {
		g.l.Error(fmt.Errorf("`g.grabGroup` something wrong: %w", err))

		return
	}

	group.LastUpdateAt = t
	err = g.groupRepo.Update(group)

	if err != nil {
		g.l.Error(fmt.Errorf("`g.grabGroup` something wrong: %w", err))
	}

	return
}

func (g *GrabberVk) newMessages(group *entity.Group) (err error) {
	for i := 0; true; i += 100 {
		rawMessages, err := g.source.GetGroupMessages(group.Name, i)
		if err != nil {
			return err
		}

		lastMessage := g.messageRepo.Last(group.ID)

		var lastMessageAt time.Time

		if lastMessage.ID == 0 {
			lastMessageAt = group.LastUpdateAt
		} else {
			lastMessageAt = lastMessage.MessageAt
		}

		for _, rawMessage := range rawMessages.Messages {
			if saved := g.saveMessage(group, rawMessage, lastMessageAt); !saved {
				return err
			}
		}
	}

	return
}

func (g *GrabberVk) saveMessage(group *entity.Group, rawMessage entity.VkMessage, lastMessageAt time.Time) bool {
	messageAt := time.Unix(rawMessage.Date, 0)
	if !lastMessageAt.Before(messageAt) {
		return false
	}

	err := g.messageRepo.Add(group.ID, rawMessage.ID, g.source.Name(), messageAt)
	if err != nil {
		g.l.Error(fmt.Errorf("`g.saveMessage`something wrong: %w", err))

		return false
	}

	g.sendMessage(group.User.TelegramID, group.Name, rawMessage)

	return true
}

func (g *GrabberVk) sendMessage(userID uint64, groupName string, message entity.VkMessage) {
	ownerID := strconv.FormatInt(message.OwnerID, 10)
	messageID := strconv.FormatUint(message.ID, 10)
	text := "https://vk.com/" + groupName + "?w=wall" + ownerID + "_" + messageID
	g.messenger.Message(userID, text)
}
