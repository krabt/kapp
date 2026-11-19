package store

import (
	"context"

	"kapp/internal/data/model"
	"kapp/internal/data/query"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetKV(ctx context.Context, keys []string) (settingMap map[string]string, err error) {
	settingMap = make(map[string]string)

	mm := query.Use(s.db).Settings
	do := mm.WithContext(ctx)
	if len(keys) > 0 {
		do = do.Where(mm.Key.In(keys...))
	}
	items, err := do.Find()
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		settingMap[item.Key] = item.Value
	}
	return settingMap, nil
}

func (s *Store) SetKV(ctx context.Context, settingMap map[string]string) error {
	if len(settingMap) == 0 {
		return nil
	}

	items := make([]model.Settings, 0, len(settingMap))
	for key, value := range settingMap {
		items = append(items, model.Settings{Key: key, Value: value})
	}

	return s.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			UpdateAll: true,
		}).Create(&items).Error
}
