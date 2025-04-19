package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func main() {
	urlArg := flag.String("u", "", "URL to load test against")
	n := flag.Int("n", 10, "Number of requests to make")

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

	err = loadTest(testURL, *n)
	if err != nil {
		fmt.Printf("running load test: %s", err.Error())
		os.Exit(1)
	}

}

func loadTest(testURL *url.URL, n int) error {
	for range n {
		resp, err := http.Get(testURL.String())
		if err != nil {
			return fmt.Errorf("making http request: %w", err)
		}

		fmt.Printf("response code: %d\n", resp.StatusCode)
	}

	return nil
}
