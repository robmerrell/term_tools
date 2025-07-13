package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v6"
)

func projectName() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	repo, err := git.PlainOpenWithOptions(cwd, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return "", fmt.Errorf("failed to open repo: %w", err)
	}

	tree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get work tree: %w", err)
	}

	return filepath.Base(tree.Filesystem.Root()), nil
}

func branchName() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	repo, err := git.PlainOpenWithOptions(cwd, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return "", fmt.Errorf("failed to open repo: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	if head.Name().IsBranch() {
		return head.Name().Short(), nil
	}

	// detatched
	return "", errors.New("detached head, make sure you are on a branch")
}
