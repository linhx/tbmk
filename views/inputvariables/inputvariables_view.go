package inputvariables

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	common "linhx.com/tbmk/common"
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

type Model struct {
	tokens              []common.Token
	inputVariableInputs map[string]InputVariableModel // TODO dont use variable name, because it can be duplicated
	focusTokenIndex     int
	focusInputVariable  *InputVariableModel
	err                 error
	quit                bool
	windowWidth         int
	windowHeight        int
}

type Tokens struct {
	Items []common.Token
}

func InitialModel(tokens Tokens) Model {
	_inputVariableInputs := make(map[string]InputVariableModel)
	focusTokenIndex := -1
	var focusInputVariable InputVariableModel
	for i := 0; i < len(tokens.Items); i++ {
		token := tokens.Items[i]
		if token.IsVariable {
			inputVariable := InputVariable(token)
			if focusTokenIndex == -1 {
				focusTokenIndex = i
				inputVariable.Focus()
				focusInputVariable = inputVariable
			}
			_inputVariableInputs[token.Name] = inputVariable
		}
	}
	return Model{
		err:                 nil,
		tokens:              tokens.Items,
		inputVariableInputs: _inputVariableInputs,
		windowWidth:         0,
		windowHeight:        0,
		focusTokenIndex:     focusTokenIndex,
		focusInputVariable:  &focusInputVariable,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) changeFocusInput(msg tea.Msg, newIndex int) tea.Cmd {
	if newIndex < 0 {
		return nil
	}

	cmds := make([]tea.Cmd, 2)
	if m.focusTokenIndex > -1 {
		variableName := m.tokens[m.focusTokenIndex].Name
		currentVariableInput := m.inputVariableInputs[variableName]
		currentVariableInput.LoseFocus()
		currentVariableInput, cmds[0] = currentVariableInput.Update(msg)
		m.inputVariableInputs[variableName] = currentVariableInput
	}

	m.focusTokenIndex = newIndex
	token := m.tokens[m.focusTokenIndex]
	theInput := m.inputVariableInputs[token.Name]
	theInput.Focus()
	theInput, cmds[1] = theInput.Update(msg)
	m.inputVariableInputs[token.Name] = theInput
	m.focusInputVariable = &theInput

	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth, m.windowHeight = msg.Width, msg.Height
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.quit = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quit = true
			return m, tea.Quit
		case tea.KeyTab:
			for i := m.focusTokenIndex + 1; i < len(m.tokens); i++ {
				if m.tokens[i].IsVariable {
					return m, m.changeFocusInput(msg, i)
				}
			}
		case tea.KeyShiftTab:
			if m.focusTokenIndex < 1 {
				return m, nil
			}
			for i := m.focusTokenIndex - 1; i > -1; i-- {
				if m.tokens[i].IsVariable {
					return m, m.changeFocusInput(msg, i)
				}
			}
		}
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.focusInputVariable != nil {
		focusTokenInput := *m.focusInputVariable
		focusTokenInput, cmd2 := focusTokenInput.Update(msg)
		m.inputVariableInputs[m.tokens[m.focusTokenIndex].Name] = focusTokenInput
		m.focusInputVariable = &focusTokenInput
		return m, cmd2
	}

	return m, nil
}

func (m Model) View() string {
	var str string
	for _, token := range m.tokens {
		if token.IsVariable {
			str += m.inputVariableInputs[token.Name].View()
		} else {
			str += token.Value
		}
	}
	style := lipgloss.NewStyle().Width(m.windowWidth - 2) // wrap text
	return fmt.Sprintf("%s", style.Render(str))
}
