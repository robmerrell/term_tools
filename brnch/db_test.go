package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	dir := t.TempDir()

	db, err := newDB(dir, "project", "branch")
	defer db.Close()

	assert.NoError(t, err)
	assert.Equal(t, "tasks_project_branch", db.prefix)
}

func TestLoadTasksMissingKey(t *testing.T) {
	dir := t.TempDir()

	db, _ := newDB(dir, "project", "branch")
	defer db.Close()

	_, err := db.LoadTasks()
	assert.NoError(t, err)
}

func TestSavingAndLoadingTasks(t *testing.T) {
	dir := t.TempDir()

	db, _ := newDB(dir, "project", "branch")
	defer db.Close()

	err := db.SaveTasks([]*Task{
		newTask("testing"),
	})
	assert.NoError(t, err)

	tasks, err := db.LoadTasks()
	assert.NoError(t, err)

	assert.Equal(t, "testing", tasks[0].Text)
}
