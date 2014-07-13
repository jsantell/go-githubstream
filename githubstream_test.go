// Similar testing in style, inspired by
// testing in github.com/google/go-github; utilities
// taken from that test suite
package githubstream

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)
import "net/http/httptest"

type values map[string]string

var (
	mux    *http.ServeMux
	server *httptest.Server
)

func setup(ghs *GithubStream) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	url, _ := url.Parse(server.URL)
	ghs.Client.BaseURL = url
	ghs.Client.UploadURL = url
}

func teardown() {
	server.Close()
}

func TestConstructorAssignment(t *testing.T) {
	frequency := time.Hour * 5
	ghs := NewGithubStream(frequency, "mozilla", "gecko-dev", "fx-team", "abcdefgh")

	if ghs.Owner != "mozilla" {
		t.Error("`Owner` property not correctly set.")
	}
	if ghs.Repo != "gecko-dev" {
		t.Error("`Repo` property not correctly set.")
	}
	if ghs.Branch != "fx-team" {
		t.Error("`Branch` property not correctly set.")
	}
	if ghs.Token != "abcdefgh" {
		t.Error("`Token` property not correctly set.")
	}
	if ghs.Frequency != frequency {
		t.Error("`Frequency` property not correctly set.")
	}
}

func TestGithubStreamStart(t *testing.T) {
	ghs := NewGithubStream(time.Second, "mozilla", "gecko-dev", "fx-team", "abcdefgh")
	setup(ghs)
	defer teardown()

	counter := 3
	var times []time.Time

	mux.HandleFunc("/repos/mozilla/gecko-dev/commits", func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		since, err := time.Parse(time.RFC3339Nano, r.Form["since"][0])

		if r.Form["sha"][0] != "fx-team" {
			t.Error("Incorrect branch used in GitHub call.")
		}

		if err != nil {
			t.Error("Could not parse `since` property.")
		}

		times = append(times, since)

		fmt.Fprintf(w, `[{"sha": "%v"}]`, counter)
	})

	for commits := range ghs.Start() {
		if *commits[0].SHA != strconv.Itoa(counter) {
			t.Error("Unexpected response from mux server.")
		}

		counter--

		if counter == 0 {
			ghs.Stop()
			close(ghs.Stream)
		}
	}

	if !times[0].Add(time.Second).Equal(times[1]) {
		t.Error("Second call should be called 1 second after the first: %v, %v", times[0], times[1])
	}
	if !times[1].Add(time.Second).Equal(times[2]) {
		t.Error("Third call should be called 1 second after the second: %v, %v", times[1], times[2])
	}
}

func TestAccessToken(t *testing.T) {
	ghs := NewGithubStream(time.Second, "mozilla", "gecko-dev", "fx-team", "")
	setup(ghs)
	defer teardown()

	mux.HandleFunc("/repos/mozilla/gecko-dev/commits", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `[{"sha":"abcdefg"}]`)
	})

	<-ghs.Start()
	ghs.Stop()
	close(ghs.Stream)
}
