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
	"github.com/gorilla/websocket"
)

const gap = "\n\n"

type errorMsg error

type receivedMsg struct {
	Message string
}

type config struct {
	host string
	port string
}

type model struct {
	viewport    viewport.Model
	messages    []string
	textarea    textarea.Model
	senderStyle lipgloss.Style
	conn        *websocket.Conn
	sub         chan []byte
	cfg         config
	err         error
}

func listenFromConn(conn *websocket.Conn, sub chan<- []byte) tea.Cmd {
	return func() tea.Msg {
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return errorMsg(err)
			}
			sub <- data
		}
	}
}

func waitForMessage(sub <-chan []byte) tea.Cmd {
	return func() tea.Msg {
		message := <-sub
		return receivedMsg{Message: string(message)}
	}
}

func writeToConn(conn *websocket.Conn, message string) tea.Cmd {
	return func() tea.Msg {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			return errorMsg(err)
		}

		return nil
	}
}

func initialModel(conn *websocket.Conn) model {
	ta := textarea.New()

	ta.Placeholder = "Send message..."
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(1)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent("Welcome to the chat room!\nType a message and press Enter to send.")

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		viewport:    vp,
		textarea:    ta,
		messages:    []string{},
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		conn:        conn,
		sub:         make(chan []byte),
		err:         nil,
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
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
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
		m.messages = append(m.messages, m.senderStyle.Render("Sender: ")+msg.Message)
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
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
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
