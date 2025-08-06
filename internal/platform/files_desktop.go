//go:build !android

package platform

import (
	"os"
	"path/filepath"
)

func UserDataDirectory(appName string) (string, error) {
	userDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(userDir, appName)

	err = os.MkdirAll(appDir, 0700)
	if err != nil {
		return "", err
	}

	return appDir, nil
}
