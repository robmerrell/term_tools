package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectName(t *testing.T) {
	name, err := projectName()
	assert.NoError(t, err)

	assert.Equal(t, "brnch", name)
}

func TestBranchName(t *testing.T) {
	branch, err := branchName()
	assert.NoError(t, err)

	assert.NotEqual(t, "", branch)
}

func TestBranchNameNoGit(t *testing.T) {
	startingDir, err := os.Getwd()
	assert.NoError(t, err)

	dir := t.TempDir()

	os.Chdir(dir)
	defer os.Chdir(startingDir)

	_, err = branchName()
	assert.Error(t, err)
}
