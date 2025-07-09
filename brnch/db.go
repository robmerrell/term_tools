package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tidwall/buntdb"
)

type DB struct {
	db     *buntdb.DB
	prefix string
}

func newDB(path, project, branch string) (*DB, error) {
	// "/home/rob/.local/share/brnch"
	// create the data dir
	if err := os.MkdirAll(path, 0750); err != nil {
		return nil, err
	}

	db, err := buntdb.Open(filepath.Join(path, "brnch.db"))
	if err != nil {
		return nil, err
	}

	prefix := fmt.Sprintf("tasks_%s_%s", project, branch)

	// create the prefix index
	pattern := fmt.Sprintf("%s_*", prefix)
	if err := db.CreateIndex(prefix, pattern, buntdb.IndexString); err != nil {
		return nil, err
	}

	return &DB{
		db:     db,
		prefix: prefix,
	}, nil
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) SaveTasks(tasks []*Task) error {
	marshalled, err := json.Marshal(tasks)
	if err != nil {
		return err
	}

	return d.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(d.taskKey(), string(marshalled), nil)
		return err
	})
}

func (d *DB) LoadTasks() ([]*Task, error) {
	tasks := []*Task{}

	err := d.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(d.taskKey())
		if err == buntdb.ErrNotFound {
			return nil
		}
		if err != nil {
			return err
		}

		return json.Unmarshal([]byte(val), &tasks)
	})

	return tasks, err
}

func (d *DB) taskKey() string {
	return fmt.Sprintf("%s_tasks", d.prefix)
}
