package track

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetTracks(t *testing.T) {
	req, err := http.NewRequest("GET", "https://api.music.yandex.net/users/974060543/playlists/3", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "OAuth AgAAAAA6Dvv_AAG8XsakJ60fUE7Hk_E5LlmW-oE")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var tResp Response
	if err = json.Unmarshal(body, &tResp); err != nil {
		t.Fatal(err)
	}

	fmt.Println(*tResp.Result.Tracks[0].Track)
}
