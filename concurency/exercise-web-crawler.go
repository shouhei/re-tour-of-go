package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type AsyncCounter struct {
	fetched map[string]bool
	num     int
	mux     sync.Mutex
}

var ac = new(AsyncCounter)
var ch = make(chan int)

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		ac.mux.Lock()
		ac.num--
		ac.mux.Unlock()
		if ac.num == -1 {
			close(ch)
		}
		return
	}
	_, ok := ac.fetched[url]
	if ok {
		ac.mux.Lock()
		ac.num--
		ac.mux.Unlock()
		if ac.num == -1 {
			close(ch)
		}
		return
	} else {
		ac.mux.Lock()
		ac.fetched[url] = true
		ac.mux.Unlock()
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		ac.mux.Lock()
		ac.num--
		ac.mux.Unlock()
		if ac.num == -1 {
			close(ch)
		}
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		ac.mux.Lock()
		ac.num++
		ac.mux.Unlock()
		go Crawl(u, depth-1, fetcher)
	}
	ac.mux.Lock()
	ac.num--
	ac.mux.Unlock()
	if ac.num == -1 {
		close(ch)
	}
	return
}

func main() {
	ac.fetched = make(map[string]bool)
	Crawl("http://golang.org/", 4, fetcher)
	for {
		_, ok := <-ch
		if !ok {
			break
		}
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
