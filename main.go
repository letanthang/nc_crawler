package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ParseResult struct {
	Title   string
	Content string
	Author  string
	Date    time.Time
	URLs    []string
}

func main() {
	ch := make(chan string)
	resultChan := make(chan ParseResult)
	crawlResult := []ParseResult{}
	crawlUrls := map[string]bool{}

	//init url
	go func() {
		ch <- "https://vnexpress.net/du-lich/mua-giang-sinh-an-tuong-o-phan-lan-voi-du-khach-viet-4032021.html"
	}()

	// share work

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for url := range ch {
			go Worker(ctx, url, resultChan)
		}
	}()

	// receive result : assign more work
	for result := range resultChan {
		// stop criteria : 5 result
		// when stop condition - close channel (ch, resultChannel)
		// panic recover

		fmt.Println("Find one", result)
		crawlResult = append(crawlResult, result)
		if len(crawlResult) == 5 {
			cancel()
			close(resultChan)
		} else {
			for _, url := range result.URLs {
				if crawlUrls[url] == false {
					crawlUrls[url] = true
					ch <- url
				}
			}
		}

	}

	fmt.Println("Num of results:", len(crawlResult))
	bs, _ := json.MarshalIndent(crawlResult, "", "  ")
	fmt.Println("Detail:", string(bs))

}

func Worker(ctx context.Context, url string, resultChan chan ParseResult) {
	temp := make(chan ParseResult)
	go func() {
		//working: parse url to content, author, urls,...
		temp <- ParseFake(url)
	}()

	select {
	case v := <-temp:
		resultChan <- v
		return
	case <-ctx.Done():
		fmt.Println("Done signal")
		return
	}

}

var m sync.Mutex

var count int

func ParseFake(url string) ParseResult {
	var i int

	m.Lock()
	count++
	i = count
	m.Unlock()

	ms := 200 + rand.Intn(300-200)
	time.Sleep(time.Duration(ms) * time.Millisecond)

	switch i {
	case 1:
		return ParseResult{Title: "Bai viet 1", Content: "content 1", URLs: []string{"http://1", "http://2", "http://3", "http://4"}}
	case 2:
		return ParseResult{Title: "Bai viet 2", Content: "content 2", URLs: []string{"http://1", "http://5", "http://6", "http://7"}}
	case 3:
		return ParseResult{Title: "Bai viet 3", Content: "content 3", URLs: []string{"http://3", "http://5", "http://8", "http://9"}}
	case 4:
		return ParseResult{Title: "Bai viet 4", Content: "content 4", URLs: []string{"http://7", "http://11", "http://12", "http://13"}}
	case 5:
		return ParseResult{Title: "Bai viet 5", Content: "content 5", URLs: []string{"http://14", "http://15", "http://16", "http://17"}}
	case 6:
		return ParseResult{Title: "Bai viet 6", Content: "content 6", URLs: []string{"http://18", "http://21", "http://17", "http://23"}}
	case 7:
		return ParseResult{Title: "Bai viet 7", Content: "content 7", URLs: []string{"http://25", "http://19", "http://27", "http://31"}}
	case 8:
		return ParseResult{Title: "Bai viet 8", Content: "content 8", URLs: []string{"http://29", "http://33", "http://41", "http://26"}}
	default:
		return ParseResult{Title: "Bai viet 10", Content: "content 10", URLs: []string{"http://34", "http://42", "http://35", "http://26"}}
	}
}
