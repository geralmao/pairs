package config

type GameContext struct {
	Version       string
	IsFxActive    bool
	IsSoundActive bool
	OutsideWidth  int
	OutsideHeight int
	Scale         float64
	OffsetX       float64
	OffsetY       float64
	Language      string
}
