package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"time"
)

func main() {
	urlArg := flag.String("u", "", "URL to load test against")
	n := flag.Int("n", 10, "Number of requests to make")
	c := flag.Int("c", 1, "Number of concurrent requests to make")

	flag.Parse()

	if *urlArg == "" {
		fmt.Printf("url must be set using -u\n")
		os.Exit(1)
	}

	testURL, err := url.Parse(*urlArg)
	if err != nil {
		fmt.Printf("parsing url: %s", err.Error())
		os.Exit(1)
	}

	results := loadTest(testURL, *n, *c)
	fmt.Printf("Successes: %d\n", results.Successes)
	fmt.Printf("Failures: %d\n", results.Failures)
}

type Response struct {
	StatusCode      int
	Error           error
	TotalTime       time.Duration
	TimeToFirstByte time.Duration
}

type Results struct {
	Successes   int
	Failures    int
	Responses   []Response
	StatusCodes map[int]int
}

func loadTest(testURL *url.URL, n int, concurrency int) Results {
	jobs := make(chan string, n)
	results := make(chan Response, n)

	// start workers
	for range concurrency {
		go worker(jobs, results)
	}

	// send work to workers
	for range n {
		jobs <- testURL.String()
	}
	close(jobs)

	res := Results{
		StatusCodes: map[int]int{},
	}
	// collect results of work
	for range n {
		resp := <-results
		if resp.Error != nil || resp.StatusCode >= 500 {
			res.Failures++
		} else {
			res.Successes++
		}
		res.Responses = append(res.Responses, resp)
		if resp.Error == nil {
			res.StatusCodes[resp.StatusCode] += 1
		}
	}

	return res
}

func worker(jobs <-chan string, results chan<- Response) {
	for j := range jobs {
		c := http.Client{}
		req, err := http.NewRequest(http.MethodGet, j, nil)
		if err != nil {
			results <- Response{
				Error:     err,
				TotalTime: 0,
			}
			return
		}

		var start time.Time
		var timeToFirstByte time.Duration
		trace := &httptrace.ClientTrace{
			GotFirstResponseByte: func() {
				timeToFirstByte = time.Since(start)
			},
		}

		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

		start = time.Now()
		resp, err := c.Do(req)
		duration := time.Since(start)
		if err != nil {
			results <- Response{
				TotalTime: duration,
				Error:     err,
			}
			return
		}
		results <- Response{
			TotalTime:       duration,
			StatusCode:      resp.StatusCode,
			TimeToFirstByte: timeToFirstByte,
		}
	}
}
