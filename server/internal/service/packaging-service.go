package service

import (
	"math"
	"server/internal/model"
	"sort"
)

const largePackItemsBuffer = 50

type PackagingService interface {
	PackItems(numberOfItems int) (map[int]int, error)
}

type PackagingServiceImpl struct {
	packsService PacksService
}

func NewPackagingService(packsService PacksService) PackagingService {
	return &PackagingServiceImpl{
		packsService: packsService,
	}
}

func (service PackagingServiceImpl) PackItems(numberOfItems int) (map[int]int, error) {
	packSizes, err := service.packsService.GetPacks()
	if err != nil {
		return nil, err
	}

	if len(packSizes) == 0 {
		return nil, &model.EmptyPacksConfig{}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(packSizes)))

	finalPacks := make(map[int]int)
	largestPack := packSizes[0]

	// Buffer to use when calculating a large number of items
	// In this approach, most of the items will go to the larget pack, but we cannot put all of them.
	// Smaller packs could provide the same number of packs as the larges, but with smaller capacity
	bufferLimit := largestPack * largePackItemsBuffer
	target := numberOfItems

	if target > bufferLimit {
		remainder := target % largestPack
		target = remainder + ((bufferLimit / largestPack) * largestPack)

		bulkCount := (numberOfItems - target) / largestPack
		finalPacks[largestPack] = bulkCount
	}

	// calculate the remaining items. Bottom-up approach is used for this.
	// the approach hee is to build an array for every item till target + largestPack + 1,
	// where in the array we are keeping how many min packs it took to package this item.
	// we are also keeping another array, that is tracking what was the last pack size used for that item.
	// this array is needed for reconstruction purposes.
	itemsToCheck := target + largestPack + 1
	minPacksForItem := make([]int, itemsToCheck)
	lastPackUsedForItem := make([]int, itemsToCheck)

	// Initialize the array with a value larger than any possible pack count
	for i := range minPacksForItem {
		minPacksForItem[i] = math.MaxInt32
	}

	minPacksForItem[0] = 0 // initial element, 0 packs are needed to make 0 items

	for i := 1; i < itemsToCheck; i++ {
		for _, pack := range packSizes {
			if i >= pack {
				if minPacksForItem[i-pack] != math.MaxInt32 {
					// Is adding this pack better than what we already found for 'i'?
					if minPacksForItem[i-pack]+1 < minPacksForItem[i] {
						minPacksForItem[i] = minPacksForItem[i-pack] + 1
						lastPackUsedForItem[i] = pack
					}
				}
			}
		}
	}

	// find the optional total, the smallest total >= target.
	bestTotal := -1
	for i := target; i < itemsToCheck; i++ {
		if minPacksForItem[i] != math.MaxInt32 {
			bestTotal = i
			break // Since we iterate up from Target, the first match is minimal items
		}
	}

	// build the final map,
	// Backtrack using the lastPackUsedForItem array
	if bestTotal != -1 {
		curr := bestTotal
		for curr > 0 {
			packUsed := lastPackUsedForItem[curr]
			finalPacks[packUsed]++
			curr -= packUsed
		}
	}

	return finalPacks, nil
}
