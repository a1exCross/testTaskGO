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

	wg.Add(len(urlSlice))

	for i := 0; i < len(urlSlice); i++ {
		ch <- struct{}{}

		go func(n int) {
			defer wg.Done()
			defer func() { <-ch }()

			count, err := getCountWordsInResponse(urlSlice[n])
			if err != nil {
				fmt.Printf("[%d] %s\n", n, err)
			} else {
				fmt.Printf("[%d] Count for \"%s\" = %d\n", n, urlSlice[n], count)
			}

			if count > 0 {
				mtx.Lock()

				answer += count

				mtx.Unlock()
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("\nTOTAL: %d", answer)
}

// Получение количества слов "Go" в теле запроса
func getCountWordsInResponse(url string) (int, error) {
	if url != "" {
		response, err := http.Get(url)
		if err != nil {
			return 0, err
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return 0, err
		}

		data := string(body)
		if data != "" {
			return strings.Count(data, "Go"), nil
		} else {
			return 0, errors.New("body is empty")
		}
	}

	return 0, errors.New("URL is empty")
}
