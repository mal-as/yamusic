package config

import "os"

type Config struct {
	UserID     string
	PlaylistID string
	Token      string
}

var Data Config

func init() {
	Data.UserID, _ = os.LookupEnv("USER_ID")
	Data.PlaylistID, _ = os.LookupEnv("PLAYLIST_ID")
	Data.Token, _ = os.LookupEnv("YA_TOKEN")
}
