package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/russross/blackfriday/v2"
	"io"
	"log"
	"net/http"
)

func stripHTMLTags(input string) string {
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

type FetchedBookmarkMsg struct {
	content string
}

func fetchLinkCmd(url string) tea.Cmd {
	return func() tea.Msg {
		markdown := fetchLink(url)
		html := blackfriday.Run([]byte(markdown))
		content := stripHTMLTags(string(html))
		return FetchedBookmarkMsg{content: content}
	}
}
