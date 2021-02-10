package playlist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mal-as/yamusic/internal/config"
	"github.com/mal-as/yamusic/internal/track"
)

type PlayList struct {
	client *http.Client
}

type Option func(p *PlayList)

func WithHTTPClient(c *http.Client) Option {
	return func(p *PlayList) {
		p.client = c
	}
}

func NewPlayList(opts ...Option) *PlayList {
	p := &PlayList{http.DefaultClient}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *PlayList) GetTracks() ([]*track.Track, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.music.yandex.net/users/%s/playlists/%s", config.Data.UserID, config.Data.PlaylistID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", config.Data.Token)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tResp track.Response
	if err = json.Unmarshal(body, &tResp); err != nil {
		return nil, err
	}

	return tResp.Result.Tracks, nil
}
