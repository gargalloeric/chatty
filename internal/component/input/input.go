package input

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textarea textarea.Model
}

func New() Model {
	ta := textarea.New()

	ta.Placeholder = "Send message ..."
	ta.Focus()

	ta.Prompt = "> "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(1)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)

	return Model{
		textarea: ta,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var taCmd tea.Cmd

	m.textarea, taCmd = m.textarea.Update(msg)

	return m, taCmd
}

func (m Model) View() string {
	return m.textarea.View()
}

func (m *Model) Reset() {
	m.textarea.Reset()
}

func (m *Model) Value() string {
	return m.textarea.Value()
}

func (m *Model) SetWidth(w int) {
	m.textarea.SetWidth(w)
}

func (m *Model) Height() int {
	return m.textarea.Height()
}
