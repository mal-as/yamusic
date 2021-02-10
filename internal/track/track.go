package track

import (
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/mal-as/yamusic/internal/config"
)

type Result struct {
	Tracks []*Track `json:"tracks"`
}

type Response struct {
	Result *Result `json:"result"`
}

type Track struct {
	Track *TrackInternal `json:"track"`
}

type TrackInternal struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Artists []*Artist `json:"artists"`
}

type Artist struct {
	Name string `json:"name"`
}

type DownloadInfoJSON struct {
	Result []struct {
		DownloadInfoUrl string `json:"downloadInfoUrl"`
		BitrateInKbps   int    `json:"bitrateInKbps"`
	} `json:"result"`
}

type DownloadInfoXML struct {
	XMLName xml.Name `xml:"download-info"`
	Host    string   `xml:"host"`
	Path    string   `xml:"path"`
	TS      string   `xml:"ts"`
	S       string   `xml:"s"`
}

func (t *Track) FullName() string {
	artSlice := make([]string, len(t.Track.Artists))

	for i := 0; i < len(t.Track.Artists); i++ {
		artSlice[i] = t.Track.Artists[i].Name
	}

	artists := strings.Join(artSlice, ",")

	return artists + " - " + t.Track.Title + ".mp3"
}

func (t *Track) Download(dir string) error {
	fmt.Println("начинаем получать информацию о загрузке")
	di, err := t.downloadInfo()
	if err != nil {
		return err
	}
	fmt.Println("информация о загрузке получена")

	fmt.Println("начинаем получать url для загрузки")
	url := di.downloadUrl()
	fmt.Println("url получен")

	client := http.DefaultClient
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", config.Data.Token)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("начинаем запись на диск")
	if err := t.saveToDisk(dir, data); err != nil {
		return err
	}
	fmt.Println("запись завершена")

	return nil
}

func (t *Track) saveToDisk(dir string, data []byte) error {
	fileName := path.Join(dir, t.FullName())
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Write(data)
	if err != nil {
		return err
	}

	if n < len(data) {
		return fmt.Errorf("writed (%d) less then expected (%d)", n, len(data))
	}

	if err = file.Sync(); err != nil {
		return err
	}

	return nil
}

func (d *DownloadInfoXML) downloadUrl() string {
	sign := fmt.Sprintf("%x", md5.Sum([]byte("XGRlBW9FXlekgbPrRHuSiA"+d.Path[1:]+d.S)))
	return fmt.Sprintf("https://%s/get-mp3/%s/%s%s", d.Host, sign, d.TS, d.Path)
}

func (t *Track) downloadInfo() (DownloadInfoXML, error) {
	url, err := t.downloadInfoURL()
	if err != nil {
		return DownloadInfoXML{}, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return DownloadInfoXML{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DownloadInfoXML{}, err
	}

	var di DownloadInfoXML
	if err := xml.Unmarshal(body, &di); err != nil {
		return DownloadInfoXML{}, err
	}

	return di, nil
}

func (t *Track) downloadInfoURL() (string, error) {
	cliet := http.DefaultClient
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.music.yandex.net/tracks/%s/download-info", t.Track.ID), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", config.Data.Token)

	resp, err := cliet.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res DownloadInfoJSON
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	for _, info := range res.Result {
		if info.BitrateInKbps == 320 {
			return info.DownloadInfoUrl, nil
		}
	}

	return "", fmt.Errorf("no download url for 320 kbps")
}
