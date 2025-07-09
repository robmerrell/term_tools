package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	tint "github.com/lrstanley/bubbletint"
)

type Task struct {
	Key     string `json:"key"`
	Text    string `json:"text"`
	Checked bool   `json:"checked"`
	// just store if this is a level 1 or level 2 task. Since we only allow 2 levels
	// this is easier than having a list of subitem structs. Moving around fully grouped
	// tasks is a little more work, but not significantly so.
	Level int `json:"level"`
}

// create a new task
func newTask(text string) *Task {
	return &Task{Text: text, Checked: false, Level: 1}
}

func (t *Task) render(m *model, selected bool) string {
	// task
	taskColor := tint.Fg()
	if t.Checked {
		taskColor = tint.BrightBlack()
	}
	if selected {
		taskColor = tint.Cyan()
	}

	task := lipgloss.NewStyle().
		Width(m.width - 6).
		Foreground(taskColor).
		Strikethrough(t.Checked).
		Render(t.Text)

	// indent for word wrapping
	task = strings.ReplaceAll(task, "\n", "\n"+strings.Repeat(" ", t.Level*2))

	return fmt.Sprintf("%s%s\n", t.renderCheckbox(), task)
}

func (t *Task) renderCheckbox() string {
	checkColor := tint.Blue()
	checkChar := "●"
	if t.Checked {
		checkChar = "✓"
		checkColor = tint.Green()
	} else {
		if t.Level == 1 {
			checkChar = "●"
		} else {
			checkChar = "○"
		}
	}

	text := fmt.Sprintf("%s%s ", strings.Repeat(" ", t.Level), checkChar)
	return lipgloss.NewStyle().Foreground(checkColor).Render(text)
}
