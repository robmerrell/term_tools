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
	Order   int    `json:"order"`
	// just store if this is a level 1 or level 2 task. Since we only allow 2 levels
	// this is easier than having a list of subitem structs. Moving around fully grouped
	// tasks is a little more work, but not significantly so.
	Level int `json:"level"`
}

// create a new task
func newTask(text string, order int) *Task {
	return &Task{Text: text, Checked: false, Level: 1, Order: order}
}

func (t *Task) render(m *model, selected bool) string {
	// checkbox
	check := "● "
	color := tint.Blue()
	if t.Checked {
		check = "✓ "
		color = tint.Green()
	}
	if t.Level == 2 {
		check = "  " + check
	}
	checkOutput := lipgloss.NewStyle().
		Foreground(color).
		Render(check)

	// task
	taskColor := tint.Fg()
	if t.Checked {
		taskColor = tint.BrightBlack()
	}
	if selected {
		taskColor = tint.Cyan()
	}

	task := lipgloss.NewStyle().
		// TODO: readdress this sizing
		Width(m.width - 10).
		Foreground(taskColor).
		Strikethrough(t.Checked).
		Render(t.Text)

	// indent for word wrapping
	task = strings.ReplaceAll(task, "\n", "\n"+strings.Repeat(" ", t.Level*2))

	return fmt.Sprintf("%s%s\n", checkOutput, task)
}
