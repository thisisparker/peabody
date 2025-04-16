package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Response structure to hold the array of strings
type Response [][]string

func main() {
	// Make a web request

	urlToRequestPtr := flag.String("url", "", "URL to request")

	flag.Parse()

	if urlToRequestPtr == nil {
		fmt.Println("Error: URL is required")
		return
	}

	fmt.Println(*urlToRequestPtr)

	cdxFormattedUrl := fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s&output=json&limit=-1&fl=timestamp", *urlToRequestPtr)

	fmt.Println(cdxFormattedUrl)

	resp, err := http.Get(cdxFormattedUrl)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response status:", resp.StatusCode)
		return
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Skip the header (first row) and get the string from the second row
	var timestampString string
	if len(response) > 1 && len(response[1]) > 0 {
		timestampString = response[1][0]
	} else {
		fmt.Println("Response doesn't have the expected structure")
	}

	fullArchiveUrl := "https://web.archive.org/web/" + timestampString + "id_/" + *urlToRequestPtr

	lastSlashIndex := strings.LastIndex(*urlToRequestPtr, "/")
	if lastSlashIndex == -1 {
		fmt.Println("Error: URL doesn't contain a slash")
		return
	}
	filename := (*urlToRequestPtr)[lastSlashIndex+1:]

	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer out.Close()

	// Get the data
	resp, err = http.Get(fullArchiveUrl)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error copying body to file:", err)
		return
	}
	fmt.Println(filename)
}
