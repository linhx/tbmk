package bookmark

import (
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	color "github.com/gookit/color"
	bookmark "linhx.com/tbmk/bookmark"
)

const MAX_DISPLAY_ITEMS int = 6

type (
	errMsg error
)
type Model struct {
	query        string
	queryInput   textinput.Model
	err          error
	bmk          bookmark.Bookmark
	matches      []bookmark.MatchedItem
	cursor       int
	firstIndex   int
	lastIndex    int
	SelectedItem bookmark.MatchedItem
	quit         bool
}

func InitialModel(bmk bookmark.Bookmark, query string) Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Prompt = "Search: "
	ti.SetValue(query)

	var matches []bookmark.MatchedItem
	if len(query) > 0 {
		matches, _ = bmk.Search(query)
	}

	return Model{
		query:      query,
		queryInput: ti,
		err:        nil,
		bmk:        bmk,
		cursor:     0,
		firstIndex: 0,
		lastIndex:  min(MAX_DISPLAY_ITEMS-1, len(matches)-1),
		matches:    matches,
	}
}

func (m *Model) jumpCursor(goDown bool) {
	if goDown {
		if m.cursor < len(m.matches)-1 {
			m.cursor++
		}
	} else {
		if m.cursor > 0 {
			m.cursor--
		}
	}
	m.setDisplayItemRange(goDown)
}

func (m *Model) resetCursor() {
	m.cursor = 0
	m.setDisplayItemRange(false)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func (m *Model) setDisplayItemRange(goDown bool) {
	if m.cursor == 0 {
		m.firstIndex = 0
		m.lastIndex = min(MAX_DISPLAY_ITEMS-1, len(m.matches)-1)
	} else {
		if goDown {
			if m.cursor > m.lastIndex {
				m.lastIndex = m.cursor
				m.firstIndex = max(0, m.lastIndex-MAX_DISPLAY_ITEMS+1)
			}
		} else {
			if m.cursor < m.firstIndex {
				m.firstIndex = m.cursor
				m.lastIndex = min(m.firstIndex+MAX_DISPLAY_ITEMS-1, len(m.matches)-1)
			}
		}
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func contains(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return true
		}
	}
	return false
}

var (
	highlightStyle       = color.Yellow
	selectedStyle        = color.BgGray
	matchedSelectedStyle = color.New(color.Yellow, color.BgGray)
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var err errMsg
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.jumpCursor(false)
		case tea.KeyDown:
			m.jumpCursor(true)
		case tea.KeyRunes:
			m.resetCursor()
		case tea.KeyBackspace:
			m.resetCursor()
		case tea.KeyDelete:
			m.resetCursor()
		case tea.KeyEnter:
			if len(m.matches) > 0 {
				m.SelectedItem = m.matches[m.cursor]
			}
			m.quit = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quit = true
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.queryInput, cmd = m.queryInput.Update(msg)
	if m.query != m.queryInput.Value() {
		m.query = m.queryInput.Value()
		m.matches, err = m.bmk.Search(m.query)
		m.resetCursor()
	}
	if err != nil {
		m.err = err
		return m, nil
	}
	return m, cmd
}

func (m Model) View() string {
	if m.quit {
		return ""
	}
	var matchesContent = strconv.Itoa(len(m.matches)) + " item(s)\n\n"
	if len(m.matches) > 0 {
		// TODO refactor check if empty MatchedIndexes then don't need to format each char
		for i := m.firstIndex; i <= m.lastIndex; i++ {
			match := m.matches[i]
			isSelected := m.cursor == i
			var line = ""
			if isSelected {
				line += selectedStyle.Render(strconv.Itoa(i+1) + ". ")
			} else {
				line += strconv.Itoa(i+1) + ". "
			}
			// format title
			_matchTitle := match.MatchTitle
			for j := 0; j < len(match.Title); j++ {
				if isSelected {
					if contains(j, _matchTitle.MatchedIndexes) {
						line += matchedSelectedStyle.Render(string(match.Title[j]))
					} else {
						line += selectedStyle.Render(string(match.Title[j]))
					}
				} else {
					if contains(j, _matchTitle.MatchedIndexes) {
						line += highlightStyle.Render(string(match.Title[j]))
					} else {
						line += string(match.Title[j])
					}
				}
			}
			// break line between tile and command
			if isSelected {
				line += selectedStyle.Render(":\n > ")
			} else {
				line += ":\n > "
			}
			// format command
			_matchCommand := match.MatchCommand
			for j := 0; j < len(match.Command); j++ {
				if isSelected {
					if contains(j, _matchCommand.MatchedIndexes) {
						line += matchedSelectedStyle.Render(string(match.Command[j]))
					} else {
						line += selectedStyle.Render(string(match.Command[j]))
					}
				} else {
					if contains(j, _matchCommand.MatchedIndexes) {
						line += highlightStyle.Render(string(match.Command[j]))
					} else {
						line += string(match.Command[j])
					}
				}
			}
			matchesContent += line + "\n"
		}
	}
	return m.queryInput.View() + "\n" + matchesContent
}
