package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mal-as/yamusic/internal/playlist"
)

func main() {
	playlist := playlist.NewPlayList(playlist.WithHTTPClient(http.DefaultClient))

	fmt.Println("начинаем получать трэки")
	tracks, err := playlist.GetTracks()
	if err != nil {
		log.Fatal(err)
	}

	if len(tracks) == 0 {
		log.Fatal("нет треков")
	}
	fmt.Println("получено", len(tracks), "трэков")

	t := tracks[0]

	fmt.Println("начинаем получать файл")
	if err = t.Download("."); err != nil {
		log.Fatal(err)
	}
	fmt.Println("файл загружен")
}
