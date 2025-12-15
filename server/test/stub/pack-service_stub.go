package stub

type PacksServiceStub struct {
	Sizes []int
	Error error
}

func (p PacksServiceStub) GetPacks() ([]int, error) {
	return p.Sizes, p.Error
}

func (p PacksServiceStub) SyncPacks(packs []int) error {
	if p.Error != nil {
		return p.Error
	}

	p.Sizes = packs
	return nil
}
