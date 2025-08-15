package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gargalloeric/chatty/internal/chat"
	"github.com/gargalloeric/chatty/internal/component/view"
	"github.com/gorilla/websocket"
)

const gap = "\n\n"

type errorMsg error

type config struct {
	host string
	port string
}

type model struct {
	view     view.Model
	textarea textarea.Model

	conn *websocket.Conn
	sub  chan chat.Message
	cfg  config
	err  error
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

	ta.KeyMap.InsertNewline.SetEnabled(false)

	view := view.New()

	return model{
		textarea: ta,
		view:     view,

		conn: conn,
		sub:  make(chan chat.Message),
		err:  nil,
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
		taCmd   tea.Cmd
		viewCmd tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.view, viewCmd = m.view.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.view.SetWidth(msg.Width)
		m.textarea.SetWidth(msg.Width)
		viewHeight := msg.Height - m.textarea.Height() - lipgloss.Height(gap)
		m.view.SetHeight(viewHeight)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			message := m.textarea.Value()
			m.textarea.Reset()
			return m, tea.Batch(taCmd, viewCmd, writeToConn(m.conn, message))
		}
	case view.TextMsg:
		return m, tea.Batch(taCmd, viewCmd, waitForMessage(m.sub))
	case errorMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(taCmd, viewCmd)
}

func (m model) View() string {
	content := fmt.Sprintf("%s%s%s", m.view.View(), gap, m.textarea.View())

	return content
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
