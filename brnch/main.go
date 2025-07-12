package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
)

// db connection
var db *DB

// errors handled in the main update loop
var lastErr error

// UI messages
type tasksSavedMsg bool // doesn't matter the type the payload isn't needed
type dbErrorMsg error
type tasksMsg []*Task

type inputMode int

const (
	// mode to move around the tasks
	normalMode inputMode = iota
	// mode to insert a new task
	insertMode
	// mode to update a an existing task
	updateMode
)

// Bubbletea app model
type model struct {
	tasks     []*Task
	taskinput textinput.Model
	cursor    int
	width     int
	height    int
	mode      inputMode
}

// Init loads the tasks from the database and sends a message to be consumed in Update
func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		tasks, err := db.LoadTasks()
		if err != nil {
			return dbErrorMsg(err)
		}

		return tasksMsg(tasks)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	stateChange := false

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.taskinput.Width = msg.Width

	case tasksMsg:
		m.tasks = msg

	case dbErrorMsg:
		lastErr = fmt.Errorf("DB Error: %w", msg)
		return m, tea.Quit

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
					cmds = append(cmds, m.saveTasksEvent())
				}

			// move the task up
			case "K":
				if m.cursor > 0 {
					m.tasks[m.cursor], m.tasks[m.cursor-1] = m.tasks[m.cursor-1], m.tasks[m.cursor]
					m.cursor--
					cmds = append(cmds, m.saveTasksEvent())
				}

			// bump the task level left
			case "H":
				m.tasks[m.cursor].Level = 1
				cmds = append(cmds, m.saveTasksEvent())

			// bump the task level right
			case "L":
				m.tasks[m.cursor].Level = 2
				cmds = append(cmds, m.saveTasksEvent())

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
					cmds = append(cmds, m.saveTasksEvent())
				}

			// toggle selection
			case " ":
				m.tasks[m.cursor].Checked = !m.tasks[m.cursor].Checked
				cmds = append(cmds, m.saveTasksEvent())
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
					task := newTask(m.taskinput.Value())
					m.tasks = append(m.tasks, task)
				} else {
					m.cursor++
					task := newTask(m.taskinput.Value())
					m.tasks = slices.Insert(m.tasks, m.cursor, task)
				}
				m.taskinput.Blur()
				m.taskinput.Reset()
				cmds = append(cmds, m.saveTasksEvent())
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

func (m *model) saveTasksEvent() tea.Cmd {
	return func() tea.Msg {
		if err := db.SaveTasks(m.tasks); err != nil {
			return dbErrorMsg(err)
		}

		return tasksSavedMsg(true)
	}
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	project, err := projectName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	branch, err := branchName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err = newDB(filepath.Join(homeDir, ".local", "share", "brnch"), project, branch)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()

	tint.NewDefaultRegistry()
	tint.SetTint(tint.TintTokyoNightStorm)

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("init error: %v", err)
		os.Exit(1)
	}

	if lastErr != nil {
		fmt.Println(lastErr)
	}
}
