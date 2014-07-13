package githubstream

import (
	"os"
	"time"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

var githubToken string = os.Getenv("FX_DEVTOOLS_BOT_GITHUB_TOKEN")

const REPO_OWNER = "mozilla"
const REPO_NAME = "gecko-dev"
const BRANCH = "master"

type GithubStream struct {
	Stream    chan []github.RepositoryCommit
	Client    *github.Client
	Ticker    *time.Ticker
	Frequency time.Duration
	Owner     string
	Repo      string
	Branch    string
	Token     string
}

func NewGithubStream(frequency time.Duration, owner string, repo string, branch string, token string) GithubStream {
	ghs := GithubStream{Frequency: frequency, Owner: owner, Repo: repo, Branch: branch, Token: token}
	ghs.Stream = make(chan []github.RepositoryCommit)
	ghs.Ticker = time.NewTicker(frequency)

	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: githubToken},
	}

	ghs.Client = github.NewClient(t.Client())

	return ghs
}

func (ghs GithubStream) Start() chan []github.RepositoryCommit {
	since := time.Now().Local().Add(-ghs.Frequency)

	go fetch(ghs, since)

	for _ = range ghs.Ticker.C {
		since = time.Now().Local().Add(-ghs.Frequency)
		go fetch(ghs, since)
	}

	return ghs.Stream
}

func (ghs GithubStream) Stop() {
	ghs.Ticker.Stop()
}

func fetch(ghs GithubStream, since time.Time) {

	opts := &github.CommitsListOptions{SHA: ghs.Branch, Since: since}
	commits, _, err := ghs.Client.Repositories.ListCommits(ghs.Owner, ghs.Repo, opts)

	if err == nil {
		ghs.Stream <- commits
	}
}