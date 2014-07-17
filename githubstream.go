package githubstream

import (
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

const VERSION = "0.2.0"

type GithubStream struct {
	Stream    chan []github.RepositoryCommit
	Client    *github.Client
	Ticker    *time.Ticker
	Frequency time.Duration
	Since     time.Duration
	Owner     string
	Repo      string
	Branch    string
	Token     string
}

func NewGithubStream(frequency time.Duration, since time.Duration, owner string, repo string, branch string, token string) *GithubStream {
	ghs := GithubStream{Frequency: frequency, Since: since, Owner: owner, Repo: repo, Branch: branch, Token: token}
	ghs.Stream = make(chan []github.RepositoryCommit)
	ghs.Ticker = time.NewTicker(frequency)

	if token != "" {
		t := &oauth.Transport{
			Token: &oauth.Token{AccessToken: ghs.Token},
		}
		ghs.Client = github.NewClient(t.Client())
	} else {
		ghs.Client = github.NewClient(nil)
	}

	return &ghs
}

func (ghs *GithubStream) Start() chan []github.RepositoryCommit {
	since := time.Now().Local().Add(-ghs.Since)

	go fetch(ghs, since)

	go func() {
		for _ = range ghs.Ticker.C {
			since = time.Now().Local().Add(-ghs.Since)
			fetch(ghs, since)
		}
	}()

	return ghs.Stream
}

func (ghs *GithubStream) Stop() {
	ghs.Ticker.Stop()
}

func fetch(ghs *GithubStream, since time.Time) {

	opts := &github.CommitsListOptions{SHA: ghs.Branch, Since: since}
	commits, _, err := ghs.Client.Repositories.ListCommits(ghs.Owner, ghs.Repo, opts)

	if err == nil {
		ghs.Stream <- commits
	}
}
