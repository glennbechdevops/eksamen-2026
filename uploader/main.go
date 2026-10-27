package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <url> <image-file>\n", os.Args[0])
		os.Exit(2)
	}
	url := os.Args[1]
	file := os.Args[2]

	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("open file: %v", err)
	}
	defer f.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image", filepath.Base(file))
	if err != nil {
		log.Fatalf("create form file: %v", err)
	}
	if _, err := io.Copy(part, f); err != nil {
		log.Fatalf("copy file: %v", err)
	}
	if err := writer.Close(); err != nil {
		log.Fatalf("close writer: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, &body)
	if err != nil {
		log.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("do request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Fprintf(os.Stderr, "status: %s\n", resp.Status)
	if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
		log.Fatalf("read response: %v", err)
	}
	fmt.Println()
}
