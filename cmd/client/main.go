package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gargalloeric/chatty/internal/chat"
	"github.com/gorilla/websocket"
)

const gap = "\n\n"

type errorMsg error

type receivedMsg chat.Message

type config struct {
	host string
	port string
}

type model struct {
	messages []string
	viewport viewport.Model
	textarea textarea.Model

	conn *websocket.Conn
	sub  chan chat.Message
	cfg  config
	err  error

	titleStyle  lipgloss.Style
	lineStyle   lipgloss.Style
	borderType  lipgloss.Border
	borderStyle lipgloss.Style
	senderStyle lipgloss.Style
}

func initialModel(conn *websocket.Conn) model {
	ta := textarea.New()

	ta.Placeholder = "Send message..."
	ta.Focus()

	ta.Prompt = "> "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(1)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent("Welcome to the chat room!\nType a message and press Enter to send.")

	ta.KeyMap.InsertNewline.SetEnabled(false)

	borderColor := lipgloss.Color("#A259EA")

	lineSyle := lipgloss.NewStyle().Foreground(borderColor)

	borderType := lipgloss.NormalBorder()

	senderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("5"))

	titleStyle := lipgloss.NewStyle().
		Background(borderColor).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Bold(true)

	borderStyle := lipgloss.NewStyle().
		Border(borderType).
		BorderForeground(borderColor)

	return model{
		messages: []string{},
		viewport: vp,
		textarea: ta,

		conn: conn,
		sub:  make(chan chat.Message),
		err:  nil,

		senderStyle: senderStyle,
		lineStyle:   lineSyle,
		titleStyle:  titleStyle,
		borderType:  borderType,
		borderStyle: borderStyle,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		listenFromConn(m.conn, m.sub),
		waitForMessage(m.sub),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		borderWidth := m.borderType.GetLeftSize() + m.borderType.GetRightSize()
		contentWidth := msg.Width - borderWidth
		m.viewport.Width = contentWidth
		m.textarea.SetWidth(contentWidth)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)

		if len(m.messages) > 0 {
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			message := m.textarea.Value()
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+message)
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			m.viewport.GotoBottom()
			return m, tea.Batch(taCmd, vpCmd, writeToConn(m.conn, message))
		}
	case receivedMsg:
		m.messages = append(m.messages, m.senderStyle.Render("Anonymous: ")+msg.Text)
		m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		m.textarea.Reset()
		m.viewport.GotoBottom()
		return m, tea.Batch(taCmd, vpCmd, waitForMessage(m.sub))
	case errorMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(taCmd, vpCmd)
}

func (m model) View() string {
	content := fmt.Sprintf("%s%s%s", m.viewport.View(), gap, m.textarea.View())

	// Render the content box without the top border
	box := m.borderStyle.BorderTop(false).Render(content)
	width := lipgloss.Width(box)

	title := m.titleStyle.Render("Test Room")

	// Build the top border line: ┌ + title + repeated dashes + ┐
	leftCorner := m.lineStyle.Render(m.borderType.TopLeft)
	rightCorner := m.lineStyle.Render(m.borderType.TopRight)
	fillWidth := width - lipgloss.Width(title) - lipgloss.Width(leftCorner) - lipgloss.Width(rightCorner)
	fill := m.lineStyle.Render(strings.Repeat(m.borderType.Top, fillWidth))
	topBorderLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftCorner,
		title,
		fill,
		rightCorner,
	)

	// Join the top border and box
	return lipgloss.JoinVertical(lipgloss.Left, topBorderLine, box)
}

func main() {
	var cfg config

	flag.StringVar(&cfg.host, "addr", "127.0.0.1", "address of the remote server")
	flag.StringVar(&cfg.port, "port", "3000", "remote server port")
	flag.Parse()

	path := fmt.Sprintf("ws://%s/v1/ws", net.JoinHostPort(cfg.host, cfg.port))

	conn, _, err := websocket.DefaultDialer.Dial(path, nil)
	if err != nil {
		fmt.Printf("There has been an error: %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(conn), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There has been an error: %v", err)
		os.Exit(1)
	}
}
