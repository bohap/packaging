package model

type EmptyPacksConfig struct {
}

func (e *EmptyPacksConfig) Error() string {
	return "no packs configured"
}
