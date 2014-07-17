go-githubstream 
================ 

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](http://godoc.org/github.com/jsantell/go-githubstream) [![Build Status](http://img.shields.io/travis/jsantell/go-githubstream.svg?style=flat)](https://travis-ci.org/jsantell/go-githubstream)

Fetch commits from a GitHub repository periodically. Uses [go-github](https://github.com/google/go-github) under the hood.


## Installation

```
$ go get github.com/jsantell/go-githubstream
```

## Documentation

Document can be found on [GoDoc](http://godoc.org/github.com/jsantell/go-githubstream)

## Example

```go
package main

import "fmt"
import "time"
import "github.com/jsantell/go-githubstream"

var TOKEN string = os.Getenv("GITHUB_ACCESS_TOKEN")

func main() {
  ghs := githubstream.NewGithubStream(time.Hour, time.Hour * 10, "jsantell", "go-githubstream", "master", TOKEN)

  // This fetches the github repo `jsantell/go-githubstream` once an hour,
  // fetching commits from that point to 10 hours prior,
  // and prints the commits
  for commits := range ghs.Start() {
    fmt.Println(commits)
  }
}
```

## License

MIT, Copyright (c) 2014 Jordan Santell
