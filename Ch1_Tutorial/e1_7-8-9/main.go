// Fetch prints the content found at each specified URL.
package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		// add https:// if it doesn't have one.
		if !strings.HasPrefix(url, "https://") || !strings.HasPrefix(url, "http://") {
			url = "https://" + url
		}
		// Fetch the url
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("fetch: %v\n", err)
		}

		// Print status code
		println(resp.Status)
		// Copy body to the os.Stdout
		if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
			log.Fatal(err)
		}

	}
}
