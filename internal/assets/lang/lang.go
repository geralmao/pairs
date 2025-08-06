package lang

import (
	"embed"
	"slices"
)

//go:embed *.json
var LanguagesFS embed.FS

type LanguageData struct {
	Menu           string `json:"menu"`
	Play           string `json:"play"`
	Settings       string `json:"settings"`
	Exit           string `json:"exit"`
	Sounds         string `json:"sounds"`
	Music          string `json:"music"`
	Mute           string `json:"mute"`
	Volume         string `json:"volume"`
	Accept         string `json:"accept"`
	Cancel         string `json:"cancel"`
	PressToPlay    string `json:"pressToPlay"`
	BackMenu       string `json:"backMenu"`
	GameOver       string `json:"gameOver"`
	YourScore      string `json:"yourScore"`
	YouWin         string `json:"youWin"`
	WallOfFame     string `json:"wallOfFame"`
	MatchType2     string `json:"matchType2"`
	MatchType3     string `json:"matchType3"`
	MatchType4     string `json:"matchType4"`
	MatchType5     string `json:"matchType5"`
	MatchType6     string `json:"matchType6"`
	MatchType7     string `json:"matchType7"`
	MatchType8     string `json:"matchType8"`
	MatchType9     string `json:"matchType9"`
	Score          string `json:"score"`
	Time           string `json:"time"`
	LevelCompleted string `json:"levelCompleted"`
	TimeBonus      string `json:"timeBonus"`
	Level          string `json:"level"`
	Goal           string `json:"goal"`
}

const LangDefault string = "en"

var langsSupported []string = []string{"es", "en"}

func IsLangIdSupported(langId string) bool {
	langIdFound := slices.ContainsFunc(langsSupported, func(langIdSupported string) bool {
		return langId == langIdSupported
	})
	return langIdFound
}
