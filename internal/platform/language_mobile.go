//go:build android

package platform

import (
	"sync"
)

func GetSystemLanguage() string {
	return androidLanguage
}

func SetLanguage(language string) {
	languageOnce.Do(func() {
		androidLanguage = language
	})
}

var languageOnce sync.Once
var androidLanguage string
