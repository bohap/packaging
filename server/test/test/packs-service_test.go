package test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"server/internal/model"
	"server/internal/service"
	"server/test/stub"
	"testing"
)

func TestGetAll_Successful(t *testing.T) {
	// given
	packs := []model.Pack{{Size: 1}, {Size: 2}}
	sizes := []int{1, 2}
	repository := stub.PacksRepositoryStub{Packs: packs}
	packsService := service.NewPacksService(repository)

	// when
	result, err := packsService.GetPacks()

	// then
	assert.Nil(t, err)
	assert.Equal(t, sizes, result)
}

func TestGetAll_Error(t *testing.T) {
	// given
	repoError := errors.New("repo error")
	repository := stub.PacksRepositoryStub{Error: repoError}
	packsService := service.NewPacksService(repository)

	// when
	result, err := packsService.GetPacks()

	// then
	assert.Nil(t, result)
	assert.Equal(t, repoError, err)
}

func TestSyncPacks_Successful(t *testing.T) {
	// given
	repository := stub.PacksRepositoryStub{}
	packsService := service.NewPacksService(repository)
	sizes := []int{1, 2}

	// when
	err := packsService.SyncPacks(sizes)

	// then
	assert.Nil(t, err)
}

func TestSyncPacks_Error(t *testing.T) {
	// given
	repoError := errors.New("repo error")
	repository := stub.PacksRepositoryStub{Error: repoError}
	packsService := service.NewPacksService(repository)

	// when
	err := packsService.SyncPacks([]int{1, 2, 3})

	// then
	assert.Equal(t, repoError, err)
}
