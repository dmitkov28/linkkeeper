package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/russross/blackfriday/v2"
)

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
		content, _ := htmltomarkdown.ConvertString(string(html))
		return FetchedBookmarkMsg{content: content}
	}
}
