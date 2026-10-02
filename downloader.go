package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// client is used for every request in this file. http.DefaultClient has no
// timeout, so a stalled connection would park the calling goroutine (and, in
// probe, the wg.Wait() goroutine) forever.
var client = &http.Client{Timeout: 15 * time.Second}

// downloader fetches the problem ZIP and reports whether it was saved, so the
// caller can skip extraction instead of reporting a second, misleading error.
func downloader(id string) bool {
	// The problem ID we want to download.

	// Resolve which language variant actually exists before building the URL.
	// probe returns "" when neither _en nor _ca serves a real ZIP, and we must
	// not fall through to a made-up URL in that case.
	lang := probe(id)
	if lang == "" {
		fmt.Println("no downloadable zip for:", id)
		return false
	}

	// Build the Jutge.org ZIP URL for that problem. The slash before "zip" is
	// required: without it the URL is .../<id>_enzip, which jutge.org answers
	// with 200 OK and an HTML page -- a soft 404 that would be saved as a file.
	url := baseURL + id + lang + "/zip"

	// Make a GET request.
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("request failed:", err)
		return false
	}
	// Always close the response body, even if we return early later.
	defer resp.Body.Close()

	// If the server didn't return 200 OK, stop.
	if resp.StatusCode != http.StatusOK {
		fmt.Println("bad status:", resp.Status)
		return false
	}

	// Status alone proves nothing on this host: unknown paths also return 200.
	// Only a real ZIP body means we got the file we asked for.
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/zip") {
		fmt.Println("not a zip (content-type:", ct, ")")
		return false
	}

	// Create (or overwrite) the local file where we'll save the ZIP.
	out, err := os.Create(id + ".zip")
	if err != nil {
		fmt.Println("create file failed:", err)
		return false
	}
	// Close the file when main exits.
	defer out.Close()

	// Copy the response body (the ZIP bytes) into the file.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("copy failed:", err)
		return false
	}

	// fmt.Println("downloaded", id+".zip")
	return true
}

func probe(id string) string {
	// Buffered so a worker never has to wait for the receiver to start
	// draining; combined with the client timeout, wg.Wait() below is bounded.
	suffixes := []string{"_en", "_ca"}
	ch := make(chan string, len(suffixes))

	var wg sync.WaitGroup
	for _, s := range suffixes {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			// Probe the ZIP endpoint itself, not the problem page. The page URL
			// (.../problems/<id>_en) answers 200 text/html for every id, even
			// ones that don't exist, so it can never distinguish anything.
			resp, err := client.Head(baseURL + id + s + "/zip")
			if err != nil {
				return
			}
			resp.Body.Close()
			// Require an actual ZIP body: jutge.org returns 200 OK for
			// nonexistent problems too, so the status code alone would always
			// look like a hit and this function would just return whichever
			// goroutine finished first.
			ct := resp.Header.Get("Content-Type")
			if resp.StatusCode == http.StatusOK && strings.HasPrefix(ct, "application/zip") {
				ch <- s
			}
		}(s)
	}

	// Close the channel once every goroutine is done, then drain it fully.
	// Taking only the first value would make the answer depend on which HEAD
	// happened to return first, so the extracted directory would flip between
	// <id>_en and <id>_ca from run to run. Collect all hits, then choose below.
	go func() { wg.Wait(); close(ch) }()
	hits := map[string]bool{}
	for s := range ch {
		hits[s] = true
	}

	// Walk suffixes in preference order instead of arrival order. Preference is
	// English, falling back to Catalan; swap the slice at the top of the
	// function to reverse it.
	for _, s := range suffixes {
		if hits[s] {
			return s
		}
	}
	// Neither variant exists: callers must treat this as "not downloadable"
	// rather than concatenating it into a URL.
	return ""
}
