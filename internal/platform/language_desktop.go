//go:build !android

package platform

import (
	"os"
	"strings"
)

// -----------------------------------------------------------------------------
// Funciones sistema
// -----------------------------------------------------------------------------
func GetSystemLanguage() string {
	langIdDetected := ""
	envVars := []string{"LC_ALL", "LC_MESSAGES", "LANG"}
	for _, key := range envVars {
		if value := os.Getenv(key); value != "" {
			// LANG suele tener formato como "es_ES.UTF-8"
			langIdDetected = strings.Split(value, ".")[0]
			langIdDetected = strings.Split(langIdDetected, "_")[0]
			break
		}
	}

	return langIdDetected
}
