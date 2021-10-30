package repo

import (
	"strconv"
	"strings"
	"time"

	"github.com/jokius/news-telegram-bot/internal/entity"
	"github.com/jokius/news-telegram-bot/pkg/postgres"
)

type UserRepo struct {
	db *postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (u UserRepo) AddGroupByURL(id, sourceName, url string) (err error) {
	user, err := u.findOrCreateUser(id)
	if err != nil {
		return
	}

	s := strings.Split(url, "/")
	groupName := s[len(s)-1]

	var group entity.Group

	u.db.Query.
		Where(&entity.Group{UserID: user.ID, SourceName: sourceName, GroupName: groupName}).
		First(&group)

	if group.ID == 0 {
		t := time.Now()
		group = entity.Group{UserID: user.ID, SourceName: sourceName, GroupName: groupName, LastUpdateAt: t,
			CreatedAt: t, UpdatedAt: t}
		err = u.db.Query.Create(&group).Error
	}

	return
}

func (u UserRepo) UpdateStartDate(id string, date time.Time) (err error) {
	user, err := u.findOrCreateUser(id)
	if err != nil {
		return
	}

	return u.db.Query.
		Model(&entity.Group{}).
		Where(&entity.Group{UserID: user.ID}).
		Updates(entity.Group{LastUpdateAt: date}).
		Error
}

func (u UserRepo) RemoveGroup(id, sourceName, url string) (err error) {
	user, err := u.findOrCreateUser(id)
	if err != nil {
		return
	}

	s := strings.Split(url, "/")
	groupName := s[len(s)-1]

	var group entity.Group
	err = u.db.Query.
		Where(&entity.Group{UserID: user.ID, SourceName: sourceName, GroupName: groupName}).
		Delete(&group).
		Error

	return
}

func (u UserRepo) Groups(id string) (groups []entity.Group, err error) {
	user, err := u.findOrCreateUser(id)
	if err != nil {
		return
	}

	err = u.db.Query.Where(&entity.Group{UserID: user.ID}).Find(&groups).Error

	return
}

func (u UserRepo) findOrCreateUser(rawID string) (user entity.User, err error) {
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil {
		return
	}

	u.db.Query.Where(&entity.User{TelegramID: id}).First(&user)

	if user.ID == 0 {
		t := time.Now()
		user.TelegramID = id
		user.CreatedAt = t
		user.UpdatedAt = t
		err = u.db.Query.Create(&user).Error
	}

	return
}
