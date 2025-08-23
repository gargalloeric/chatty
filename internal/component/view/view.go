package view

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var senderStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("5"))

type TextMsg struct {
	Sender string
	Text   string
}

type MetadataMsg struct {
	Room      string
	UserCount int
}

type Model struct {
	messages []string
	viewport viewport.Model

	roomTitle string
	userCount int

	borderType  lipgloss.Border
	borderStyle lipgloss.Style
	titleStyle  lipgloss.Style
	lineStyle   lipgloss.Style
}

func New() Model {
	vp := viewport.New(30, 5)

	vp.SetContent("Welcome to the chat room!\nType a message and press Enter to send.")

	borderColor := lipgloss.Color("#A259EA")
	borderType := lipgloss.NormalBorder()
	titleStyle := lipgloss.NewStyle().
		Background(borderColor).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Bold(true)
	borderStyle := lipgloss.NewStyle().
		Border(borderType).
		BorderForeground(borderColor)
	lineStyle := lipgloss.NewStyle().
		Foreground(borderColor)

	return Model{
		messages:    []string{},
		viewport:    vp,
		borderType:  borderType,
		borderStyle: borderStyle,
		titleStyle:  titleStyle,
		lineStyle:   lineStyle,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var vpCmd tea.Cmd

	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case TextMsg:
		m.messages = append(m.messages, senderStyle.Render(fmt.Sprintf("%s: ", msg.Sender))+msg.Text)
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.viewport.GotoBottom()
	case MetadataMsg:
		m.roomTitle = msg.Room
		m.userCount = msg.UserCount
	}

	return m, vpCmd
}

func (m Model) View() string {
	content := m.viewport.View()

	// Render the content box without the top border
	box := m.borderStyle.BorderTop(false).Render(content)
	width := lipgloss.Width(box)

	title := m.titleStyle.Render(m.roomTitle)
	userCount := m.titleStyle.Render(fmt.Sprintf("Users %d", m.userCount))
	separator := m.lineStyle.Render(m.borderType.Top)

	// Build the top border line: ┌ + title + dash + user count + repeated dashes + ┐
	leftCorner := m.lineStyle.Render(m.borderType.TopLeft)
	rightCorner := m.lineStyle.Render(m.borderType.TopRight)
	elementsWidth := lipgloss.Width(title) + lipgloss.Width(separator) + lipgloss.Width(userCount) + lipgloss.Width(leftCorner) + lipgloss.Width(rightCorner)
	fillWidth := width - elementsWidth
	fill := m.lineStyle.Render(strings.Repeat(m.borderType.Top, fillWidth))
	topBorderLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftCorner,
		title,
		separator,
		userCount,
		fill,
		rightCorner,
	)

	// Join the top border and box
	return lipgloss.JoinVertical(lipgloss.Left, topBorderLine, box)
}

func (m *Model) SetWidth(w int) {
	borderWidth := m.borderType.GetLeftSize() + m.borderType.GetRightSize()
	m.viewport.Width = w - borderWidth

	if len(m.messages) > 0 {
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
	}
	m.viewport.GotoBottom()
}

func (m *Model) SetHeight(h int) {
	borderHeight := m.borderType.GetTopSize() + m.borderType.GetBottomSize()
	m.viewport.Height = h - borderHeight
}
