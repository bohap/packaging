package stub

import (
	"server/internal/model"
)

type PacksRepositoryStub struct {
	Packs []model.Pack
	Error error
}

func (p PacksRepositoryStub) FindAll() ([]model.Pack, error) {
	return p.Packs, p.Error
}

func (p PacksRepositoryStub) SyncPacks(packs []int) error {
	return p.Error
}
