package service

import (
	"log"
	"server/internal/repository"
)

type PacksService interface {
	GetPacks() ([]int, error)
	SyncPacks(packs []int) error
}

type PacksServiceImpl struct {
	repository repository.PacksRepository
}

func NewPacksService(repository repository.PacksRepository) PacksService {
	return &PacksServiceImpl{
		repository: repository,
	}
}

func (service PacksServiceImpl) GetPacks() ([]int, error) {
	packs, err := service.repository.FindAll()
	if err != nil {
		return nil, err
	}

	sizes := make([]int, len(packs))
	for i, p := range packs {
		sizes[i] = p.Size
	}

	return sizes, nil
}

func (service PacksServiceImpl) SyncPacks(packs []int) error {
	log.Printf("Syncing packs: %v", packs)

	if err := service.repository.SyncPacks(packs); err != nil {
		log.Printf("Error syncing packs: %v", err)
		return err
	}

	return nil
}
