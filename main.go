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
	TotalTimes  struct {
		Max  time.Duration
		Min  time.Duration
		Mean time.Duration
		All  []time.Duration
	}
	TimeToFirstBytes struct {
		Max  time.Duration
		Min  time.Duration
		Mean time.Duration
		All  []time.Duration
	}
}

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
	printResults(results)
}

func loadTest(testURL *url.URL, n int, concurrency int) Results {
	jobCh := make(chan string, n)
	resultCh := make(chan Response, n)

	// start workers
	for range concurrency {
		go worker(jobCh, resultCh)
	}

	// send work to workers
	for range n {
		jobCh <- testURL.String()
	}
	close(jobCh)

	results := Results{
		StatusCodes: map[int]int{},
	}
	// collect results of work
	for range n {
		resp := <-resultCh
		collectMetrics(&results, resp)
	}

	return results
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

func collectMetrics(results *Results, response Response) {
	if response.Error != nil {
		results.Failures++
		return
	}
	results.StatusCodes[response.StatusCode] += 1
	results.TotalTimes.All = append(results.TotalTimes.All, response.TotalTime)
	results.TimeToFirstBytes.All = append(results.TimeToFirstBytes.All, response.TimeToFirstByte)
	if response.StatusCode >= 500 {
		results.Failures++
	} else {
		results.Successes++
	}

	results.TotalTimes.Max = max(results.TotalTimes.Max, response.TotalTime)
	if results.TotalTimes.Min == 0 {
		results.TotalTimes.Min = response.TotalTime
	}
	results.TotalTimes.Min = min(results.TotalTimes.Min, response.TotalTime)
	results.TimeToFirstBytes.Max = max(results.TimeToFirstBytes.Max, response.TimeToFirstByte)
	if results.TimeToFirstBytes.Min == 0 {
		results.TimeToFirstBytes.Min = response.TimeToFirstByte
	}
	results.TimeToFirstBytes.Min = min(results.TimeToFirstBytes.Min, response.TimeToFirstByte)
}

func printResults(results Results) {
	fmt.Printf("Results\n")
	fmt.Printf("\tSuccess (2xx) %d\n", results.Successes)
	fmt.Printf("\tFailures (5xx) %d\n", results.Failures)

	var totalTime int64
	for _, d := range results.TotalTimes.All {
		totalTime += int64(d)
	}
	meanTotalTime := time.Duration(totalTime / int64(len(results.TotalTimes.All)))
	fmt.Printf("Total Request Time (max, min, mean): %f %f %f\n", results.TotalTimes.Max.Seconds(), results.TotalTimes.Min.Seconds(), meanTotalTime.Seconds())
	var totalTimeToFirstByte int64
	for _, d := range results.TimeToFirstBytes.All {
		totalTimeToFirstByte += int64(d)
	}
	meanTimeToFirstByte := time.Duration(totalTime / int64(len(results.TotalTimes.All)))
	fmt.Printf("Time To First Byte (max, min, mean): %f %f %f\n", results.TimeToFirstBytes.Max.Seconds(), results.TimeToFirstBytes.Min.Seconds(), meanTimeToFirstByte.Seconds())
}
