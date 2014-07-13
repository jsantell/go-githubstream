go-githubstream
===============

Fetch commits from a GitHub repository periodically. Uses [go-github](https://github.com/google/go-github) under the hood.

[Documentation](http://godoc.org/github.com/jsantell/go-githubstream)

## Installation

```
$ go get github.com/jsantell/go-githubstream
```

## Example

```go
package main

import "fmt"
import "time"
import "github.com/jsantell/go-githubstream"

var TOKEN string = os.Getenv("GITHUB_ACCESS_TOKEN")

func main() {
  ghs := githubstream.NewGithubStream(time.Hour, "jsantell", "go-githubstream", "master", TOKEN)

  // This fetches the github repo `jsantell/go-githubstream` once an hour
  // and prints the commits
  for commits := range ghs.Start() {
    fmt.Println(commits)
  }
}
```

## License

MIT, Copyright (c) 2014 Jordan Santell
