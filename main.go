package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/russross/blackfriday/v2"
)

type ViewMode int

const (
	ListView ViewMode = iota
	AddMode
	BookmarkView
)

type model struct {
	mode         ViewMode
	addInput     textinput.Model
	bookmarks    list.Model
	bookmarkView viewport.Model
	// err        error
}

type bookmark struct {
	url string
}

func (b bookmark) Title() string       { return b.url }
func (b bookmark) Description() string { return "" }
func (b bookmark) FilterValue() string { return b.url }

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter bookmark URL"
	ti.CharLimit = 156
	ti.Width = 50

	items := []list.Item{
		bookmark{url: "https://example.com"},
	}

	bookmarkList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	bookmarkList.Title = "Bookmarks"

	bookmarkView := viewport.New(200, 200)
	bookmarkView.SetContent("\n\nThis is the bookmark view")

	return model{
		mode:         ListView,
		addInput:     ti,
		bookmarks:    bookmarkList,
		bookmarkView: bookmarkView,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "a":
			if m.mode == ListView {
				m.mode = AddMode
				m.addInput.Focus()
				return m, textinput.Blink
			}

		case "l":
			if m.mode == AddMode {
				m.mode = ListView
				m.addInput.Blur()
				return m, nil
			}

		case "enter":
			if m.mode == AddMode && m.addInput.Value() != "" {
				newBookmark := bookmark{url: m.addInput.Value()}
				m.bookmarks.InsertItem(0, newBookmark)

				m.addInput.SetValue("")
				m.addInput.Blur()
				m.mode = ListView
				return m, nil
			}

			if m.mode == ListView {
				selectedItem, ok := m.bookmarks.SelectedItem().(bookmark)
				if ok {
					markdown := fetchLink(selectedItem.url)

					html := blackfriday.Run([]byte(markdown))

					content := stripHTMLTags(string(html))
					m.bookmarkView.SetContent(content + "\n\n" + selectedItem.url)
				}
				m.mode = BookmarkView
				return m, nil
			}

		case "backspace":
			if m.mode == BookmarkView {
				m.mode = ListView
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		h, v := lipgloss.NewStyle().Margin(2, 2).GetFrameSize()
		m.bookmarks.SetSize(msg.Width-h, msg.Height-v)
		m.bookmarkView.Width = msg.Width - h
		m.bookmarkView.Height = msg.Height - v
	}

	switch m.mode {
	case AddMode:
		newInput, cmd := m.addInput.Update(msg)
		m.addInput = newInput
		cmds = append(cmds, cmd)

	case ListView:
		newList, cmd := m.bookmarks.Update(msg)
		m.bookmarks = newList
		cmds = append(cmds, cmd)

	case BookmarkView:

	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.mode {
	case AddMode:
		return fmt.Sprintf(
			"Add New Bookmark\n\n%s\n\nPress Enter to save, L to view list",
			m.addInput.View(),
		)

	case ListView:
		return fmt.Sprintf(
			"%s\n\nPress A to add new bookmark, Q to quit",
			m.bookmarks.View(),
		)

	case BookmarkView:
		return m.bookmarkView.View()

	default:
		return "Unknown view mode"
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
