package variableinputs

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

var titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("148")).Foreground(lipgloss.Color("236"))

// TODO refactor Model structure
type Model struct {
	tokens                  []common.Token
	focusTokenIndex         int
	variableInputs          map[int]VariableInputModel
	focusVariableInputIndex int
	focusVariableInput      *VariableInputModel
	err                     error
	Quit                    bool
	Cancel                  bool
	windowWidth             int
	windowHeight            int
}

func InitialModel(command string, windowWidth int, windowHeight int) Model {
	tokens := common.TokensParser(command)
	variableInputs := make(map[int]VariableInputModel)
	focusTokenIndex := -1
	focusVariableInputIndex := -1
	var focusVariableInput VariableInputModel
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if token.IsVariable {
			variableInput := VariableInput(token)
			if focusTokenIndex == -1 {
				focusTokenIndex = i
				variableInput.Focus()
				focusVariableInput = variableInput
				focusVariableInputIndex = 0
			}
			variableInputs[token.Id] = variableInput
		}
	}
	return Model{
		err:                     nil,
		tokens:                  tokens,
		variableInputs:          variableInputs,
		focusVariableInputIndex: focusVariableInputIndex,
		focusTokenIndex:         focusTokenIndex,
		focusVariableInput:      &focusVariableInput,
		windowWidth:             windowWidth,
		windowHeight:            windowHeight,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) GetValue() string {
	str := ""
	for _, token := range m.tokens {
		str += token.Value
	}
	return str
}

func (m *Model) updateTokenValueById(id int, value string) {
	for i := 0; i < len(m.tokens); i++ {
		if m.tokens[i].Id == id {
			m.tokens[i].Value = value
		}
	}
}

func (m *Model) changeFocusInput(msg tea.Msg, newIndex int) tea.Cmd {
	if newIndex < 0 {
		return nil
	}

	cmds := make([]tea.Cmd, 2)
	// blur current variable input
	if m.focusVariableInput != nil {
		tokenId := m.tokens[m.focusTokenIndex].Id
		currentVariableInput := *m.focusVariableInput
		currentVariableInput.LoseFocus()
		currentVariableInput, cmds[0] = currentVariableInput.Update(msg)
		m.variableInputs[tokenId] = currentVariableInput // override
	}

	// focus to the new one
	m.focusTokenIndex = newIndex
	tokenId := m.tokens[m.focusTokenIndex].Id
	newVariableInput := m.variableInputs[tokenId]
	newVariableInput.Focus()
	newVariableInput, cmds[1] = newVariableInput.Update(msg)
	m.variableInputs[tokenId] = newVariableInput // override
	m.focusVariableInput = &newVariableInput

	return tea.Batch(cmds...)
}

func (m Model) hasVariable() bool {
	for _, e := range m.tokens {
		if e.IsVariable {
			return true
		}
	}
	return false
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.hasVariable() {
		m.Quit = true
		return m, tea.Quit
	}
	m.Quit = false
	m.Cancel = false
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth, m.windowHeight = msg.Width, msg.Height
		return m, tea.ClearScreen
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.Quit = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Cancel = true
			return m, tea.Quit
		case tea.KeyTab:
			for i := m.focusTokenIndex + 1; i < len(m.tokens); i++ {
				if m.tokens[i].IsVariable {
					m.focusVariableInputIndex++
					return m, m.changeFocusInput(msg, i)
				}
			}
		case tea.KeyShiftTab:
			if m.focusTokenIndex < 1 {
				return m, nil
			}
			for i := m.focusTokenIndex - 1; i > -1; i-- {
				if m.tokens[i].IsVariable {
					m.focusVariableInputIndex--
					return m, m.changeFocusInput(msg, i)
				}
			}
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	// update focused input
	if m.focusVariableInput != nil {
		focusTokenInput := *m.focusVariableInput
		focusTokenInput, cmd2 := focusTokenInput.Update(msg)
		m.variableInputs[focusTokenInput.token.Id] = focusTokenInput
		m.focusVariableInput = &focusTokenInput
		m.updateTokenValueById(focusTokenInput.token.Id, focusTokenInput.input.Value())
		return m, cmd2
	}

	return m, nil
}

func removeFirstNLines(input string, n int) string {
	if n < 1 {
		return input
	}
	lineCount := 0
	for i := 0; i < len(input); i++ {
		if input[i] == '\n' {
			lineCount++
			if lineCount == n {
				return input[i+1:]
			}
		}
	}
	return ""
}

func (m Model) getLineNumberOfVariableToken(token common.Token) int {
	var content string
	for _, _token := range m.tokens {
		if _token == token {
			content += "."
			break
		}
		content += _token.Value
	}

	renderStr := lipgloss.NewStyle().Width(m.windowWidth - 2).Render(content)
	count := 0
	for _, char := range renderStr {
		if char == '\n' {
			count++
		}
	}
	return count
}

func (m Model) View() string {
	var str string
	for _, token := range m.tokens {
		if token.IsVariable {
			str += m.variableInputs[token.Id].View()
		} else {
			str += token.Value
		}
	}
	style := lipgloss.NewStyle().Width(m.windowWidth - 2) // wrap text
	str = style.Render(str)
	var totalVariables = len(m.variableInputs)
	numOfVariables := titleStyle.Render(fmt.Sprintf("Variables: %d/%d", m.focusVariableInputIndex+1, totalVariables))
	return numOfVariables + "\n" + style.MaxHeight(5).Render(removeFirstNLines(str, m.getLineNumberOfVariableToken(m.tokens[m.focusTokenIndex])-1))
}
