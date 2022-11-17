package views_save

import (
	"strings"

	common "linhx.com/tbmk/common"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	bookmark "linhx.com/tbmk/bookmark"
)

type Model struct {
	focusIndex          int
	inputs              []textinput.Model
	cursorMode          textinput.CursorMode
	err                 error
	bmk                 bookmark.Bookmark
	quit                bool
	confirmOverrideMode bool
	confirmOverrideMsg  string
}

func (m *Model) GetItem() (string, string) {
	return m.inputs[0].Value(), m.inputs[1].Value()
}

func InitialModel(bmk bookmark.Bookmark, command string) Model {
	m := Model{
		inputs:              make([]textinput.Model, 2),
		confirmOverrideMode: false,
		bmk:                 bmk,
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()

		switch i {
		case 0:
			t.CharLimit = 100
			t.Prompt = "Title: "
			t.Focus()
		case 1:
			t.CharLimit = 500
			t.Prompt = "Command: "
			t.SetValue(command)
		}

		m.inputs[i] = t
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

/**
 * return true if want to stay
 */
func (m *Model) save(override bool) (bool, error) {
	title, command := m.GetItem()
	if len(title) > 0 && len(command) > 0 {
		_, err := m.bmk.Save(title, command, override)
		_, ok := err.(*common.DuplicateBmkiError)
		if ok {
			m.confirmOverrideMode = true
			m.confirmOverrideMsg = err.Error()
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (m Model) updateOverrideMode(msg tea.Msg) (tea.Model, tea.Cmd) {
	var err error
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			_, err = m.save(true)
			if err != nil {
				m.err = err
			}
			m.confirmOverrideMode = false
			m.confirmOverrideMsg = ""
			m.quit = true
			return m, tea.Quit
		case "n", "ctrl+c":
			m.confirmOverrideMode = false
			m.confirmOverrideMsg = ""
		}
	case error:
		m.err = msg
		return m, nil
	}

	return m, nil
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) updateInputView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quit = true
			return m, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := range m.inputs {
				cmds[i] = m.inputs[i].SetCursorMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		case "enter":
			shouldStay, err := m.save(false)
			if err != nil {
				m.err = err
			}
			if shouldStay {
				return m, nil
			} else {
				m.quit = true
				return m, tea.Quit
			}
		// Set focus to next input
		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.confirmOverrideMode {
		return m.updateOverrideMode(msg)
	} else {
		return m.updateInputView(msg)
	}
}

func (m Model) ConfirmOverrideView() string {
	return m.confirmOverrideMsg + ". Do you want to update the existed item (y/n)"
}

func (m Model) View() string {
	if m.confirmOverrideMode {
		return m.ConfirmOverrideView()
	}
	if m.quit {
		return ""
	}

	var b strings.Builder
	b.WriteString("Commands bookmark\n")
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}
