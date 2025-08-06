package config

import "image/color"

const WindowWidth int = 420
const WindowHeight int = 760

const Dt float32 = 1.0 / 60.0

const OffsetX float32 = 9
const OffsetY float32 = 45
const CardWidth float32 = 100
const CardHeight float32 = 100
const CardPoints uint = 20
const ScoreBonusPoints uint = 35

var BackgroundColorApplication color.NRGBA = color.NRGBA{0x7a, 0x5b, 0x9c, 0xff}
