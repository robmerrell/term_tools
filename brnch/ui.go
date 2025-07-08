package main

import (
	"slices"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
)

type inputMode int

const (
	normalMode inputMode = iota
	insertMode
)

type model struct {
	tasks     []*Task
	taskinput textinput.Model
	cursor    int
	width     int
	height    int
	mode      inputMode
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	stateChange := false

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.taskinput.Width = msg.Width

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if m.mode == normalMode {
			switch msg.String() {
			// move the cursor up
			case "k":
				if m.cursor > 0 {
					m.cursor--
				}

			// move the cursor down
			case "j":
				if m.cursor < len(m.tasks)-1 {
					m.cursor++
				}

			// move the task down
			case "J":
				if m.cursor < len(m.tasks)-1 {
					m.tasks[m.cursor], m.tasks[m.cursor+1] = m.tasks[m.cursor+1], m.tasks[m.cursor]
					m.cursor++
				}

			// move the task up
			case "K":
				if m.cursor > 0 {
					m.tasks[m.cursor], m.tasks[m.cursor-1] = m.tasks[m.cursor-1], m.tasks[m.cursor]
					m.cursor--
				}

			// bump the task level left
			case "H":
				m.tasks[m.cursor].Level = 1

			// bump the task level right
			case "L":
				m.tasks[m.cursor].Level = 2

			// focus task creation
			case "i":
				m.mode = insertMode
				cmd := m.taskinput.Focus()
				cmds = append(cmds, cmd)
				stateChange = true

			// delete
			case "d":
				if len(m.tasks) > 0 {
					m.tasks = slices.Delete(m.tasks, m.cursor, m.cursor+1)
					if m.cursor > len(m.tasks)-1 {
						m.cursor--
					}
				}

			// toggle selection
			case " ":
				m.tasks[m.cursor].Checked = !m.tasks[m.cursor].Checked
			}
		} else {
			switch msg.Type {
			case tea.KeyEsc:
				m.mode = normalMode
				m.taskinput.Blur()
				m.taskinput.Reset()

			case tea.KeyEnter:
				m.mode = normalMode
				if m.cursor == len(m.tasks) {
					task := newTask(m.taskinput.Value(), m.cursor)
					m.tasks = append(m.tasks, task)
				} else {
					m.cursor++
					task := newTask(m.taskinput.Value(), m.cursor)
					m.tasks = slices.Insert(m.tasks, m.cursor, task)
				}
				m.taskinput.Blur()
				m.taskinput.Reset()
			}
		}
	}

	var cmd tea.Cmd
	m.taskinput, cmd = m.taskinput.Update(msg)
	if stateChange {
		m.taskinput.Reset()
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	header := lipgloss.NewStyle().
		Background(tint.Bg()).
		Foreground(tint.Blue()).
		Align(lipgloss.Right).
		Width(m.width).
		Render("term_tools :: master")

	textInput := lipgloss.NewStyle().
		Width(m.width).
		Render(m.taskinput.View())

	footer := lipgloss.NewStyle().
		Background(tint.Bg()).
		Foreground(tint.Fg()).
		Align(lipgloss.Left).
		Width(m.width).
		Render("i: new task  u: update task  space: toggle  d: delete  HJKL: move task")

	inner := ""
	for i, task := range m.tasks {
		inner += task.render(&m, i == m.cursor)
	}

	contentMargin := 2
	content := lipgloss.NewStyle().
		Width(m.width - contentMargin*2).
		Height(m.height - lipgloss.Height(header) - lipgloss.Height(textInput) - lipgloss.Height(footer) - contentMargin*2).
		Align(lipgloss.Left).
		Margin(2).
		Render(inner)

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer, textInput)
}

func initialModel() model {
	input := textinput.New()
	input.Placeholder = "Task"

	return model{
		mode:      normalMode,
		tasks:     []*Task{},
		cursor:    0,
		taskinput: input,
	}
}
