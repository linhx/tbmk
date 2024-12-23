/*
 * This file is part of tbmk.
 *
 * tbmk is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * tbmk is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with tbmk.  If not, see <https://www.gnu.org/licenses/>.
 *
 * Copyright (C) 2024 Nguyen Dinh Linh
 */

package bookmark

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	color "github.com/gookit/color"
	bookmark "linhx.com/tbmk/bookmark"
	common "linhx.com/tbmk/common"
	variableinputs "linhx.com/tbmk/views/variableinputs"
)

const MAX_DISPLAY_ITEMS int = 6
const ITEM_HEIGHT = 2
const HEADER_HEIGHT = 3

/*
 * Sample content of 1 item:
 * Search:
 * N item(s)
 * ---------------
 * 1. item 1
 *  > command 1
 */
const MIN_WINDOW_HEIGHT = HEADER_HEIGHT + ITEM_HEIGHT // 1 search input, 1 number of item(s), 1 hr, 2 for one item

type (
	errMsg error
)

type state int

const (
	initializing         state  = 0
	ready                state  = 1
	SEARCH_MODE          string = "SEARCH"
	DELETE_MODE          string = "DELETE"
	INPUT_VARIABLES_MODE string = "VARIABLES"
)

var (
	topLabelStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("148")).Background(lipgloss.Color("236")).MarginRight(1)
	highlightStyle       = color.Yellow
	selectedStyle        = color.BgGray
	matchedSelectedStyle = color.New(color.Yellow, color.BgGray)
	ellipsisStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	listCommandsStyle    = lipgloss.NewStyle().Width(20).Border(lipgloss.NormalBorder(), false, true, false, false).BorderForeground(lipgloss.Color("69"))
	previewCommandStyle  = lipgloss.NewStyle().Width(20)
	previewTitleStyle    = lipgloss.NewStyle().Background(lipgloss.Color("148")).Foreground(lipgloss.Color("236"))
)

type Model struct {
	state                    state
	query                    string
	queryInput               textinput.Model
	err                      error
	bmk                      bookmark.Bookmark
	matches                  []bookmark.MatchedItem
	cursor                   int // index of selected item. index in `matches`
	numberOfDisplayableItems int
	firstIndex               int // first item index of current displayed items
	lastIndex                int // last item index of current displayed items
	SelectedItem             bookmark.MatchedItem
	OutputCommand            string
	quit                     bool
	windowWidth              int
	windowHeight             int
	activateMode             string
	variablesInputView       *variableinputs.Model
}

func InitialModel(bmk bookmark.Bookmark, query string) Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	ti.Prompt = "TBMK - Search:"
	ti.PromptStyle = topLabelStyle
	ti.SetValue(query)

	var matches []bookmark.MatchedItem
	matches, _ = bmk.Search(query)

	return Model{
		state:                    initializing,
		activateMode:             SEARCH_MODE,
		query:                    query,
		queryInput:               ti,
		err:                      nil,
		bmk:                      bmk,
		cursor:                   0,
		numberOfDisplayableItems: 0,
		firstIndex:               0,
		lastIndex:                0,
		matches:                  matches,
		windowWidth:              0,
		windowHeight:             0,
	}
}

func (m *Model) init() {
	var heightForItems = m.windowHeight - HEADER_HEIGHT
	var newNumberOfDisplayableItems = min(heightForItems/ITEM_HEIGHT, MAX_DISPLAY_ITEMS)

	m.numberOfDisplayableItems = newNumberOfDisplayableItems
	m.lastIndex = m.numberOfDisplayableItems - 1
}

func (m *Model) reNumberOfDisplayableItems() {
	var heightForItems = m.windowHeight - HEADER_HEIGHT
	m.numberOfDisplayableItems = min(heightForItems/ITEM_HEIGHT, MAX_DISPLAY_ITEMS)
}

func (m Model) getContentHeight() int {
	return m.numberOfDisplayableItems * ITEM_HEIGHT
}

func (m *Model) reCalcCursor() {
	matchesLastIndex := len(m.matches) - 1
	if m.cursor >= matchesLastIndex {
		m.cursor = matchesLastIndex
	}
	if m.lastIndex > matchesLastIndex {
		m.lastIndex = matchesLastIndex
	}
	m.firstIndex = max(m.lastIndex-m.numberOfDisplayableItems+1, 0)
}

func (m *Model) moveCursor(goDown bool) {
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
		m.lastIndex = min(m.numberOfDisplayableItems-1, len(m.matches)-1)
	} else {
		if goDown {
			if m.cursor > m.lastIndex {
				m.lastIndex = m.cursor
				m.firstIndex = max(0, m.lastIndex-m.numberOfDisplayableItems+1)
			}
		} else {
			if m.cursor < m.firstIndex {
				m.firstIndex = m.cursor
				m.lastIndex = min(m.firstIndex+m.numberOfDisplayableItems-1, len(m.matches)-1)
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

func (m Model) updateDeleteMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err errMsg
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			m.bmk.Remove(m.matches[m.cursor].Id)
			m.matches, err = m.bmk.Search(m.query)
			m.reCalcCursor()
			m.activateMode = SEARCH_MODE
		case "n", "ctrl+c":
			m.activateMode = SEARCH_MODE
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	if err != nil {
		m.err = err
		return m, nil
	}
	return m, nil
}

func (m Model) updateSearchMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var err errMsg
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp, tea.KeyShiftTab:
			m.moveCursor(false)
		case tea.KeyDown, tea.KeyTab:
			m.moveCursor(true)
		case tea.KeyRunes:
			m.resetCursor()
		case tea.KeyBackspace:
			m.resetCursor()
		case tea.KeyDelete:
			m.resetCursor()
		case tea.KeyEnter:
			if len(m.matches) > 0 {
				m.SelectedItem = m.matches[m.cursor]
				// open variables input view
				variablesInput := variableinputs.InitialModel(m.SelectedItem.Command, m.windowWidth, m.windowHeight)
				m.variablesInputView = &variablesInput
				m.activateMode = INPUT_VARIABLES_MODE

				return m.updateInputVariablesMode(nil) // avoid escalate KeyEnter to variableinputs_view
			}
			m.quit = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quit = true
			return m, tea.Quit
		case tea.KeyCtrlD:
			m.activateMode = DELETE_MODE
			return m, cmd
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

func (m Model) updateInputVariablesMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	variablesInputView := *m.variablesInputView
	variablesInputView, _ = variablesInputView.Update(msg)
	if variablesInputView.Quit {
		m.variablesInputView = nil
		m.activateMode = SEARCH_MODE
		m.quit = true
		m.OutputCommand = variablesInputView.GetValue()
		return m, tea.Quit
	}
	if variablesInputView.Cancel {
		m.activateMode = SEARCH_MODE
		m.variablesInputView = nil
	}
	m.variablesInputView = &variablesInputView
	return m, nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth, m.windowHeight = msg.Width, msg.Height
		m.state = ready
		m.reNumberOfDisplayableItems()
		m.resetCursor()
		contentHeight := m.getContentHeight()
		listCommandsStyle = listCommandsStyle.MaxHeight(contentHeight).Width((m.windowWidth / 5) * 2)
		previewCommandStyle = previewCommandStyle.MaxHeight(contentHeight).Width(m.windowWidth - listCommandsStyle.GetWidth())
	}

	switch m.activateMode {
	case INPUT_VARIABLES_MODE:
		return m.updateInputVariablesMode(msg)
	case DELETE_MODE:
		return m.updateDeleteMode(msg)
	default:
		return m.updateSearchMode(msg)
	}
}

func (m Model) DeleteView() string {
	return "Do you want to delete item " + selectedStyle.Render(m.matches[m.cursor].Command) + "? (y/n)"
}

func (m Model) View() string {
	if m.state == initializing {
		return "..."
	}
	if m.windowHeight < MIN_WINDOW_HEIGHT {
		return "Window height is not enough to display"
	}

	switch m.activateMode {
	case INPUT_VARIABLES_MODE:
		return (*m.variablesInputView).View()
	case DELETE_MODE:
		return m.DeleteView()
	}
	if m.quit {
		return ""
	}
	var matchesContent = strconv.Itoa(len(m.matches)) + " item(s)\n----------"
	var content string
	if len(m.matches) > 0 {
		// TODO refactor check if empty MatchedIndexes then don't need to format each char
		listCommandsContent := ""
		for i := m.firstIndex; i <= m.lastIndex; i++ {
			match := m.matches[i]
			isSelected := m.cursor == i
			var line = "\n"
			if i == m.firstIndex {
				line = ""
			}
			if isSelected {
				line += selectedStyle.Render(strconv.Itoa(i+1) + ". ")
			} else {
				line += strconv.Itoa(i+1) + ". "
			}
			// format title
			_matchTitle := match.MatchTitle
			truncatedTitle, isTruncatedTitle := common.TruncateWithEllipsis(match.Title, listCommandsStyle.GetWidth())
			for j := 0; j < len(truncatedTitle); j++ {
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
			if isTruncatedTitle {
				if isSelected {
					line += selectedStyle.Render(ellipsisStyle.Render(common.ELLIPSIS))
				} else {
					line += ellipsisStyle.Render(common.ELLIPSIS)
				}
			}
			// break line between tile and command
			if isSelected {
				line += selectedStyle.Render(":") + "\n" + selectedStyle.Render(" > ")
			} else {
				line += ":\n > "
			}
			// format command
			_matchCommand := match.MatchCommand
			truncatedCommand, isTruncated := common.TruncateWithEllipsis(match.Command, listCommandsStyle.GetWidth())
			for j := 0; j < len(truncatedCommand); j++ {
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
			if isTruncated {
				if isSelected {
					line += selectedStyle.Render(ellipsisStyle.Render(common.ELLIPSIS))
				} else {
					line += ellipsisStyle.Render(common.ELLIPSIS)
				}
			}
			listCommandsContent += line
		}

		selectedCommandStr := m.matches[m.cursor].Command

		content = lipgloss.JoinHorizontal(lipgloss.Top, listCommandsStyle.MaxHeight(m.getContentHeight()).Render(listCommandsContent), previewCommandStyle.Render(previewTitleStyle.Render(" Preview ")+"\n"+selectedCommandStr))
	}
	return fmt.Sprintf("%s\n%s\n%s", m.queryInput.View(), matchesContent, content)
}
