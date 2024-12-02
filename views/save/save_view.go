package views_save

import (
	common "linhx.com/tbmk/common"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bookmark "linhx.com/tbmk/bookmark"
)

const TOP_HEIGHT = 4        // 1 for app title, 1 for title input, 1 for command prompt
const MIN_WINDOW_HEIGHT = 5 // 1 for app title, 1 for title input, 2 for command input
const MAX_COMMAND_INPUT_HEIGHT = 6

var (
	topLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("148")).Background(lipgloss.Color("236")).MarginRight(1)
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type Model struct {
	focusIndex          int
	titleInput          textinput.Model
	commandInput        textarea.Model
	err                 error
	bmk                 bookmark.Bookmark
	quit                bool
	confirmOverrideMode bool
	confirmOverrideMsg  string
	windowWidth         int
	windowHeight        int
}

const (
	NUM_INPUTS = 2
)

func (m *Model) GetItem() (string, string) {
	return m.titleInput.Value(), m.commandInput.Value()
}

func InitialModel(bmk bookmark.Bookmark, command string) Model {
	titleInput := textinput.New()
	titleInput.CharLimit = 100
	titleInput.Prompt = "Title: "
	titleInput.Focus()

	commandInput := textarea.New()
	commandInput.SetValue(command)

	return Model{
		confirmOverrideMode: false,
		bmk:                 bmk,
		titleInput:          titleInput,
		commandInput:        commandInput,
	}
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
	cmds := make([]tea.Cmd, NUM_INPUTS)

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	m.titleInput, cmds[0] = m.titleInput.Update(msg)
	m.commandInput, cmds[1] = m.commandInput.Update(msg)

	return tea.Batch(cmds...)
}

func (m Model) updateInputView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quit = true
			return m, tea.Quit

		case tea.KeyCtrlS:
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
		case tea.KeyTab, tea.KeyShiftTab:
			// Cycle indexes
			if msg.Type == tea.KeyShiftTab {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > NUM_INPUTS {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = NUM_INPUTS
			}

			cmds := make([]tea.Cmd, NUM_INPUTS)
			switch m.focusIndex {
			case 0:
				cmds[0] = m.titleInput.Focus()
				m.commandInput.Blur()
				break
			case 1:
				m.titleInput.Blur()
				cmds[1] = m.commandInput.Focus()
				break
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth, m.windowHeight = msg.Width, msg.Height
		commandInputHeight := m.windowHeight - TOP_HEIGHT
		if commandInputHeight > MAX_COMMAND_INPUT_HEIGHT {
			commandInputHeight = MAX_COMMAND_INPUT_HEIGHT
		}
		m.commandInput.SetHeight(commandInputHeight)
		m.commandInput.SetWidth(m.windowWidth)
	}
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
	if m.windowHeight < MIN_WINDOW_HEIGHT {
		return "Window height is not enough to display"
	}
	if m.confirmOverrideMode {
		return m.ConfirmOverrideView()
	}
	if m.quit {
		return ""
	}

	return topLabelStyle.Render("TBMK - Save") + "\n" + helpStyle.Render("(Ctrl + S to save)") + "\n" + m.titleInput.View() + "\n" + "Command: \n" + m.commandInput.View()
}
