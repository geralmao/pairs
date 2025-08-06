//go:build js && wasm

package storage

func SaveGameData(data *GameData) error {
	return nil
}

func LoadGameData() (*GameData, error) {
	return &GameData{
		MusicActive: true,
		FxActive:    true,
		Language:    "en",
	}, nil
}
