package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

const MAX_GOROUTINES = 5

func main() {
	urlSlice := []string{"https://golang.org", "https://golang.org", "https://golang.org",
		"https://golang.org", "https://golang.org", "https://golang.org", "https://golang.org",
		"https://golang.org", "https://golang.org", "", "blablabla", "not_wrong_url"}

	var ch = make(chan struct{}, MAX_GOROUTINES)
	var wg sync.WaitGroup
	var mtx sync.Mutex
	var answer int = 0

	for i := 0; i < len(urlSlice); i++ {
		ch <- struct{}{}
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			count, err := sendRequest(urlSlice[n])
			if err != nil {
				fmt.Printf("[%d] %s\n", n, err)
			} else {
				fmt.Printf("[%d] Count for %s = %d\n", n, urlSlice[n], count)
			}

			mtx.Lock()

			if count > 0 {
				answer += count
			}

			mtx.Unlock()

			<-ch
		}(i)
	}

	wg.Wait()

	fmt.Println(answer)
}

func sendRequest(url string) (int, error) {
	if url != "" {
		res, err := http.Get(url)
		if err != nil {
			return -1, err
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return -1, err
		}

		data := string(body)
		if data != "" {
			return strings.Count(data, "Go"), nil
		} else {
			return -1, errors.New("body is empty")
		}
	}

	return -1, errors.New("URL is empty")
}
