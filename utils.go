package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func stripHTMLTags(input string) string {
	// Simple and naive way of stripping HTML tags
	// This could be improved for more complex cases
	var result string
	inTag := false
	for _, c := range input {
		switch c {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				result += string(c)
			}
		}
	}
	return result
}

func fetchLink(url string) string {
	const jinaAiUrl = "https://r.jina.ai"
	modifiedUrl := fmt.Sprintf("%s/%s", jinaAiUrl, url)
	bytes, err := http.Get(modifiedUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer bytes.Body.Close()
	res, err := io.ReadAll(bytes.Body)

	if err != nil {
		log.Fatal(err)
	}
	return string(res)
}

