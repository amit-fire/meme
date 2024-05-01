package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Result struct {
	ID  string
	Val []string
}

type Callable struct {
	ID string
}

func (m Callable) call() Result {
	val := req(m.ID, m.ID, m.ID, m.ID)
	return Result{ID: m.ID, Val: val}
}

func run(numberOfRequests int) {
	var wg sync.WaitGroup
	results := make(chan Result, 100)

	t := getTime()
	for i := 0; i < numberOfRequests; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			callable := Callable{ID: strconv.Itoa(i + 1)}
			results <- callable.call()
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Printf("id %s code %s content %s\n", result.ID, result.Val[0], result.Val[1])
	}

	fmt.Println("start time " + t)
	fmt.Println("end time " + getTime())
}

func getTime() string {
	now := time.Now()
	return now.Format("2006/01/02 15:04:05:000")
}

func req(id, lat, lon, query string) []string {
	url := fmt.Sprintf("http://localhost:7000/app/memes?lat=%s&lon=%s&query=%s", lat, lon, query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return []string{"", ""}
	}
	req.Header.Set("id", id)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return []string{"", ""}
	}
	defer resp.Body.Close()

	responseCode := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return []string{"", ""}
	}

	content := strings.TrimSpace(string(body))
	return []string{fmt.Sprintf("%d", responseCode), content}
}

func main() {
	run(100)
}
