package play

type Level struct {
	target int
	time   float32
	emojis int
	matchs int
}

func NewLevel(target int, time float32, emojis, matchs int) *Level {
	return &Level{target: target, time: time, emojis: emojis, matchs: matchs}
}

func (l Level) Target() int {
	return l.target
}

func (l Level) Time() float32 {
	return l.time
}

func (l Level) EmojisCount() int {
	return l.emojis
}

func (l Level) Matchs() int {
	return l.matchs
}
