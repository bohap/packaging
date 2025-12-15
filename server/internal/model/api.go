package model

type PacksSyncRequest struct {
	Packs []int `json:"packs" binding:"required"`
}

type ProductsPackageRequest struct {
	NumberOfItems int `json:"numberOfItems" binding:"required"`
}

type ProductPackageResponse map[int]int
