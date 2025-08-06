package storage

type GameData struct {
	MusicActive bool   `json:"musicActive"`
	FxActive    bool   `json:"fxActive"`
	Language    string `json:"language"`
}
