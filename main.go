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

	if *urlArg == ""{
		fmt.Printf("url must be set using -u\n")
		os.Exit(1)
	}

	url, err := url.Parse(*urlArg)
	if err != nil {
		fmt.Printf("parsing url: %s", err.Error())
		os.Exit(1)
	}

	for range *n {
		resp, err := http.Get(url.String())
		if err != nil {
			fmt.Printf("making http request: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("response code: %d\n", resp.StatusCode)
	}
}
