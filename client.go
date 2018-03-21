package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

var (
	client = &http.Client{Timeout: 10 * time.Second}
)

func getBody(url string) (bytes.Buffer, error) {
	var buf bytes.Buffer

	r, err := client.Get(url)
	if err != nil {
		return buf, err
	}
	defer r.Body.Close()

	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		return buf, err
	}

	return buf, nil
}

func poll() (routines Routines, err error) {
	url := fmt.Sprintf("http://%s/%s/goroutine?debug=2", *hostFlag, *endpointFlag)
	buf, err := getBody(url)
	if err != nil {
		return
	}

	return ReadRoutines(buf), nil
}
