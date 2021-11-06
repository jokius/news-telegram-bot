package repo

import (
	"time"

	"github.com/jokius/news-telegram-bot/internal/entity"
	"github.com/jokius/news-telegram-bot/pkg/postgres"
)

type GroupRepo struct {
	db *postgres.Postgres
}

func NewGroupRepo(pg *postgres.Postgres) *GroupRepo {
	return &GroupRepo{pg}
}

func (g GroupRepo) AllBySource(source string) (groups []entity.Group, err error) {
	err = g.db.Query.
		Preload("User").
		Model(&entity.Group{}).
		Where(&entity.Group{SourceName: source}).
		Find(&groups).Error

	return
}

func (g GroupRepo) Update(group *entity.Group) (err error) {
	group.UpdatedAt = time.Now().UTC()

	return g.db.Query.Save(group).Error
}
