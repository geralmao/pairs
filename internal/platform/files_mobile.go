//go:build android

package platform

import (
	"sync"
)

func UserDataDirectory(appName string) (string, error) {
	return androidDataPath, nil
}

func SetUserDataPath(dataPath string) {
	pathOnce.Do(func() {
		androidDataPath = dataPath
	})
}

var pathOnce sync.Once
var androidDataPath string
