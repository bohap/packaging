package test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"server/internal/model"
	service2 "server/internal/service"
	"server/test/stub"
	"testing"
)

func TestPackItems(t *testing.T) {
	scenarios := []struct {
		name        string
		packSizes   []int
		itemsToPack int
		expected    map[int]int
	}{
		{
			name:        "simple pack: matches smallest size",
			packSizes:   []int{100, 200, 1000},
			itemsToPack: 100,
			expected:    map[int]int{100: 1},
		},
		{
			name:        "simple pack: matches middle size",
			packSizes:   []int{100, 200, 1000},
			itemsToPack: 200,
			expected:    map[int]int{200: 1},
		},
		{
			name:        "simple pack: matches largest size",
			packSizes:   []int{100, 200, 1000},
			itemsToPack: 1000,
			expected:    map[int]int{1000: 1},
		},
		{
			name:        "simple pack: less then smallest back",
			packSizes:   []int{100, 200, 1000},
			itemsToPack: 1,
			expected:    map[int]int{100: 1},
		},
		{
			name:        "simple pack: bigger then smallest, less then next one",
			packSizes:   []int{100, 200, 1000},
			itemsToPack: 101,
			expected:    map[int]int{200: 1},
		},
		{
			name:        "simple pack: bigger then largest",
			packSizes:   []int{100, 200, 1000},
			itemsToPack: 1001,
			expected:    map[int]int{1000: 1, 100: 1},
		},
		{
			name:        "simple pack: matches smallest size",
			packSizes:   []int{100, 200, 1000},
			itemsToPack: 100,
			expected:    map[int]int{100: 1},
		},
		{
			name:        "different pack choices",
			packSizes:   []int{250, 500, 1000, 2000, 5000},
			itemsToPack: 12001,
			expected:    map[int]int{5000: 2, 2000: 1, 250: 1},
		},
		{
			name:        "large number of items",
			packSizes:   []int{250, 500, 1000, 2000, 5000},
			itemsToPack: 500001,
			expected:    map[int]int{5000: 100, 250: 1},
		},
		{
			name:        "edge case",
			packSizes:   []int{23, 31, 53},
			itemsToPack: 500000,
			expected:    map[int]int{23: 2, 31: 7, 53: 9429},
		},
	}

	for _, scenario := range scenarios {
		t.Run(
			scenario.name, func(t *testing.T) {
				// given
				packsServiceStub := stub.PacksServiceStub{Sizes: scenario.packSizes}
				service := service2.NewPackagingService(packsServiceStub)

				// when
				result, err := service.PackItems(scenario.itemsToPack)

				// then
				assert.Nil(t, err)
				assert.Equal(t, scenario.expected, result)
			},
		)
	}
}

func TestPackItems_PacksFetchError(t *testing.T) {
	// given
	serviceErr := errors.New("error")
	packsServiceStub := stub.PacksServiceStub{Error: serviceErr}
	service := service2.NewPackagingService(packsServiceStub)

	// when
	result, err := service.PackItems(1)

	// then
	assert.Equal(t, serviceErr, err)
	assert.Nil(t, result)
}

func TestPackItems_EmptyPacksList(t *testing.T) {
	// given
	packsServiceStub := stub.PacksServiceStub{Sizes: []int{}}
	service := service2.NewPackagingService(packsServiceStub)

	// when
	result, err := service.PackItems(1)

	// then
	assert.Equal(t, &model.EmptyPacksConfig{}, err)
	assert.Nil(t, result)
}
