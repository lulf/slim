/*
 * Copyright 2019, Ulf Lilleengen
 * License: Apache License 2.0 (see the file LICENSE or http://apache.org/licenses/LICENSE-2.0.html).
 */
package datastore

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func (ds SqlDatastore) Close() {
	ds.handle.Close()
}

func NewSqliteDatastore(fileName string, maxSize int) (*SqlDatastore, error) {

	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Print("Opening Database:", err)
		return nil, err
	}

	return &SqlDatastore{
		handle:  db,
		maxSize: maxSize,
	}, nil
}

func (ds SqlDatastore) Initialize() error {
	// Create initial database table
	tableCreate := `
	create table if not exists events (id integer not null primary key, insertion_time integer, creation_time integer, device_id integer, payload text);
        `

	_, err := ds.handle.Exec(tableCreate)
	if err != nil {
		log.Print("Creating Database Tables:", err)
		return err
	}
	return nil
}

func (ds SqlDatastore) InsertEvent(event *Event) error {
	tx, err := ds.handle.Begin()
	if err != nil {
		log.Print("Starting transaction:", err)
		return err
	}

	insertStmt, err := tx.Prepare("INSERT INTO events(insertion_time, creation_time, device_id, payload) values(?, ?, ?, ?)")
	if err != nil {
		log.Print("Preparing insert statement:", err)
		return err
	}
	defer insertStmt.Close()

	removeStmt, err := tx.Prepare("DELETE FROM events WHERE device_id=? AND id NOT IN (SELECT id FROM events WHERE device_id=? ORDER BY insertion_time DESC LIMIT ?)")
	if err != nil {
		log.Print("Preparing remove statement:", err)
		return err
	}
	defer removeStmt.Close()

	_, err = insertStmt.Exec(event.insertTime, event.creationTime, event.deviceId, event.payload)
	if err != nil {
		log.Print("Inserting entry:", err)
		return err
	}

	_, err = removeStmt.Exec(event.deviceId, event.deviceId, ds.maxSize)
	if err != nil {
		log.Print("Removing oldest entry:", err)
		return err
	}

	return tx.Commit()
}

func (ds SqlDatastore) ListEvents(offset int) ([]*Event, error) {
	return nil, nil
}

func (ds SqlDatastore) NumEvents() (int, error) {
	return 0, nil
}

func (ds SqlDatastore) WatchEvents(offset int, watcher Watcher) (Watch, error) {
	return nil, nil
}
