package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("argument must be provided\n")
		os.Exit(1)
	}

	url, err := url.Parse(args[0])
	if err != nil {
		fmt.Printf("parsing url: %s", err.Error())
		os.Exit(1)
	}

	resp, err := http.Get(url.String())
	if err != nil {
		fmt.Printf("making http request: %s", err.Error())
		os.Exit(1)
	}

	fmt.Printf("response code: %d", resp.StatusCode)
}
