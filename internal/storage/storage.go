//go:build !js && !wasm

package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/programatta/pairs/internal/platform"
)

// SaveGameData guarda la información del objeto GameData en el fichero
// data.json y devuelve un error en caso que se produzca.
func SaveGameData(data *GameData) error {
	//en Android nos puede devolver la ruta del tipo:
	// - /data/user/0/paquete.de.aplicacion/
	// - /storage/emulated/0/Android/data/paquete.de.aplicacion/files/
	//en desktop nos puede devolver la ruata del tipo:
	// - linux  : /home/usuario/.config/
	// - windows: C:\Users\usuario\AppData\Roaming\
	appDir, err := platform.UserDataDirectory(appName)
	if err != nil {
		return err
	}

	filePath := filepath.Join(appDir, fileData)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("\nSaveGameData: error creando el fichero [%s] en directorio [%v] - con error:[%v]", fileData, appDir, err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	return encoder.Encode(data)
}

// LoadGameData carga el fichero data.json y devuelve la información en un
// objeto GameData o error.
func LoadGameData() (*GameData, error) {
	appDir, err := platform.UserDataDirectory(appName)
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(appDir, fileData)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var gameData GameData
	err = json.NewDecoder(file).Decode(&gameData)
	return &gameData, err
}

const appName string = "pairs"
const fileData string = "data.json"
