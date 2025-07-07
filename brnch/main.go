package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

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

type listItem struct {
	task    string
	checked bool
	// just store if this is a level 1 or level 2 task. Since we only allow 2 levels
	// this is easier than having a list of subitem structs. Moving around fully grouped
	// tasks is a little more work, but not significantly so.
	level int
	order int
}

func newTask(text string, cursor int) *listItem {
	return &listItem{checked: false, level: 1, task: text, order: cursor}
}

func (l *listItem) render(m *model, selected bool) string {
	// checkbox
	check := "● "
	color := tint.Blue()
	if l.checked {
		check = "✓ "
		color = tint.Green()
	}
	if l.level == 2 {
		check = "  " + check
	}
	checkOutput := lipgloss.NewStyle().
		Foreground(color).
		Render(check)

	// task
	taskColor := tint.Fg()
	if l.checked {
		taskColor = tint.BrightBlack()
	}
	if selected {
		taskColor = tint.Cyan()
	}

	task := lipgloss.NewStyle().
		// TODO: readdress this sizing
		Width(m.width - 10).
		Foreground(taskColor).
		Strikethrough(l.checked).
		Render(l.task)

	// indent for word wrapping
	task = strings.ReplaceAll(task, "\n", "\n"+strings.Repeat(" ", l.level*2))

	return fmt.Sprintf("%s%s\n", checkOutput, task)
}

type model struct {
	list      []*listItem
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
				if m.cursor < len(m.list)-1 {
					m.cursor++
				}

			// move the task down
			case "J":
				if m.cursor < len(m.list)-1 {
					m.list[m.cursor], m.list[m.cursor+1] = m.list[m.cursor+1], m.list[m.cursor]
					m.cursor++
				}

			// move the task up
			case "K":
				if m.cursor > 0 {
					m.list[m.cursor], m.list[m.cursor-1] = m.list[m.cursor-1], m.list[m.cursor]
					m.cursor--
				}

			// bump the task level left
			case "H":
				m.list[m.cursor].level = 1

			// bump the task level right
			case "L":
				m.list[m.cursor].level = 2

			// focus task creation
			case "i":
				m.mode = insertMode
				cmd := m.taskinput.Focus()
				cmds = append(cmds, cmd)
				stateChange = true

			// delete
			case "d":
				if len(m.list) > 0 {
					m.list = slices.Delete(m.list, m.cursor, m.cursor+1)
					if m.cursor > len(m.list)-1 {
						m.cursor--
					}
				}

			// toggle selection
			case " ":
				m.list[m.cursor].checked = !m.list[m.cursor].checked
			}
		} else {
			switch msg.Type {
			case tea.KeyEsc:
				m.mode = normalMode
				m.taskinput.Blur()
				m.taskinput.Reset()

			case tea.KeyEnter:
				m.mode = normalMode
				m.cursor++
				task := newTask(m.taskinput.Value(), m.cursor)
				m.list = slices.Insert(m.list, m.cursor, task)
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
	for i, item := range m.list {
		inner += item.render(&m, i == m.cursor)
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
		mode: normalMode,
		list: []*listItem{
			{checked: false, level: 1, task: "hello"},
			{checked: false, level: 1, task: "two"},
			{checked: false, level: 2, task: "sub1 kajsd lfkjas dlfkj asldfkj alsdkfj alsdkfj laskdjf laksdjf aklsdjf laksjdf laksdjf alskdjf laskdjf laksdjf lkasjd flkasjd fkajsdf"},
			{checked: true, level: 2, task: "sub2"},
			{checked: false, level: 2, task: "sub3"},
			{checked: false, level: 1, task: "three kj asdlkfj a;sldkfj l;askdjf ljoiwer oiwuer oiwue roiwu eroiwu eroiuw eroiu weoriuw eroiuweori uwoeiru woeiruwoeiru owieru oiweurowieur"},
		},
		cursor:    0,
		taskinput: input,
	}
}

func main() {
	tint.NewDefaultRegistry()
	tint.SetTint(tint.TintTokyoNightStorm)

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("init error: %v", err)
		os.Exit(1)
	}
}
