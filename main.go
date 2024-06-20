package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const path = "khodams.txt"

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	width   int
	height  int
	name    textinput.Model
	khodams []string
	pairs   map[string]string
	table   table.Model
}

func load_khodams() []string {
	var results []string
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return results
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	return results
}

func load_model() model {
	ti := textinput.New()
	ti.Placeholder = "John Doe"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	columns := []table.Column{
		{Title: "Nama", Width: 40},
		{Title: "Khodam", Width: 40},
	}

	rows := []table.Row{}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	t.SetStyles(s)
	return model{
		name:    ti,
		khodams: load_khodams(),
		pairs:   make(map[string]string),
		table:   t,
	}
}

func (m *model) generate_khodam(name string) {
	if khodam, ok := m.pairs[name]; ok {
		m.table.SetRows(append([]table.Row{[]string{name, khodam}}, m.table.Rows()...))
		return
	}
	splitName := strings.Split(name, "")
	r, _ := utf8.DecodeRuneInString(splitName[len(splitName)-1])
	rand := len(splitName) + int(r)

	khodam := m.khodams[rand%len(m.khodams)]
	m.pairs[name] = khodam
	m.table.SetRows(append([]table.Row{[]string{name, khodam}}, m.table.Rows()...))
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			name := m.name.Value()
			if name == "" {
				return m, nil
			}
			m.generate_khodam(name)
			m.name.Reset()
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	m.name, cmd = m.name.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			"Cek Khodam",
			m.name.View(),
			baseStyle.Render(m.table.View()),
			"(esc) to quit",
		),
	)
}

func main() {
	p := tea.NewProgram(load_model())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
