package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"server/internal/model"
)

type PacksRepository interface {
	FindAll() ([]model.Pack, error)
	SyncPacks(packs []int) error
}

type PacksRepositoryImpl struct {
	db *gorm.DB
}

func (repo *PacksRepositoryImpl) FindAll() ([]model.Pack, error) {
	var packs []model.Pack
	result := repo.db.Model(&model.Pack{}).
		Order("size asc").
		Find(&packs)
	return packs, result.Error
}

func NewPacksRepository(db *gorm.DB) PacksRepository {
	return &PacksRepositoryImpl{db: db}
}

func (repo *PacksRepositoryImpl) SyncPacks(packs []int) error {
	return repo.db.Transaction(
		func(tx *gorm.DB) error {
			if len(packs) == 0 {
				return tx.Where("1 = 1").Delete(&model.Pack{}).Error
			}

			// delete old
			if err := tx.Where("size NOT IN ?", packs).Delete(&model.Pack{}).Error; err != nil {
				return err
			}

			// insert new
			var toSave []model.Pack
			for _, s := range packs {
				toSave = append(toSave, model.Pack{Size: s})
			}

			if err := tx.Clauses(
				clause.OnConflict{DoNothing: true},
			).Create(&toSave).Error; err != nil {
				return err
			}

			return nil
		},
	)
}
