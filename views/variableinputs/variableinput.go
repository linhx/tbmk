package variableinputs

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	common "linhx.com/tbmk/common"
)

type VariableInputModel struct {
	token       common.Token
	input       textinput.Model
	hasFocus    bool
	err         error
	edited      bool
	selectedAll bool
	defaultSet  bool
}

var (
	defaultStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("148")).Background(lipgloss.Color("236"))
	placeholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("23")).Background(lipgloss.Color("236"))
	hasFocusStyle    = defaultStyle.Underline(true)
	selectedAllStyle = hasFocusStyle.Underline(true).Background(lipgloss.Color("21")).Foreground(lipgloss.Color("255"))
)

func VariableInput(token common.Token) VariableInputModel {
	t := textinput.New()
	t.CharLimit = 0
	t.SetValue(token.Value)
	t.Placeholder = "{{" + token.Name + "}}" // TODO avoid duplicate code
	t.TextStyle = defaultStyle
	t.Cursor.TextStyle = defaultStyle
	t.PlaceholderStyle = placeholderStyle
	t.Prompt = ""
	return VariableInputModel{
		hasFocus:   false,
		err:        nil,
		token:      token,
		input:      t,
		edited:     false,
		defaultSet: true,
	}
}

func (m *VariableInputModel) SetFocus(focus bool) {
	if focus {
		m.Focus()
	} else {
		m.LoseFocus()
	}
}

func (m *VariableInputModel) Focus() {
	m.hasFocus = true
	m.input.Focus()
	m.input.TextStyle = selectedAllStyle
	m.input.Cursor.TextStyle = selectedAllStyle
	m.input.CursorEnd()
	m.defaultSet = true
}

func (m *VariableInputModel) LoseFocus() {
	m.hasFocus = false
	m.selectedAll = false
	m.input.CursorStart()
	m.input.Blur()
	m.input.TextStyle = defaultStyle
	m.input.Cursor.TextStyle = defaultStyle
}

func (m VariableInputModel) Update(msg tea.Msg) (VariableInputModel, tea.Cmd) {
	// TODO duplicate code
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.defaultSet && (msg.Type == tea.KeyRunes || msg.Type == tea.KeyBackspace || msg.Type == tea.KeyDelete) {
			m.input.SetValue("") // Clear the default value
			m.defaultSet = false // Mark as edited
			m.input.TextStyle = defaultStyle
			m.input.Cursor.TextStyle = defaultStyle
			m.input.CursorStart()
		} else if msg.Type == tea.KeyRight || msg.Type == tea.KeyLeft || msg.Type == tea.KeyHome || msg.Type == tea.KeyEnd || msg.Type == tea.KeyCtrlE {
			m.defaultSet = false
			m.input.TextStyle = defaultStyle
			m.input.Cursor.TextStyle = defaultStyle
		} else if msg.Type == tea.KeyCtrlA { // select all
			m.defaultSet = true
			m.input.TextStyle = selectedAllStyle
			m.input.Cursor.TextStyle = selectedAllStyle
			m.input.CursorEnd()
		}
	}

	newInput, cmd := m.input.Update(msg)
	m.input = newInput
	m.token.Value = newInput.Value()
	m.edited = true
	return m, cmd
}

func (m VariableInputModel) View() string {
	return m.input.View()
}
