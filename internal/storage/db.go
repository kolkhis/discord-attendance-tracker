package storage

import (
	"fmt"
	// "os"
	// "path/filepath"
	"database/sql"
	// _ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func Open(path string) (*DB, error) {
	fmt.Printf("Open database at path: %s\n", path)
	// os.MkdirAll(filepath.Dir(path), 0o755) // Creates the directory for the path if it doesn't exist
	return &DB{}, nil // Change this to return an actual database connection
}

func (db *DB) Close() error {
	fmt.Printf("Closing database connection\n")
	return nil // Change this to close the actual database connection
}

func (db *DB) initSchema() error {
	// Implement the database schema initialization logic here
	fmt.Printf("Initializing database schema\n")
	const schema = `
CREATE TABLE IF NOT EXISTS events (
	event_id TEXT PRIMARY KEY,
	guild_id TEXT NOT NULL,
	channel_id TEXT,
	name TEXT NOT NULL,
	entity_type INTEGER NOT NULL,
	scheduled_start_time TEXT NOT NULL,
	scheduled_end_time TEXT,
	tracking_open_time TEXT,
	tracking_close_time TEXT,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS event_subscriptions (
	event_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	subscribed_at TEXT NOT NULL,
	PRIMARY KEY (event_id, user_id)
);

CREATE TABLE IF NOT EXISTS voice_sessions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	event_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	channel_id TEXT NOT NULL,
	joined_at TEXT NOT NULL,
	left_at TEXT
);

CREATE TABLE IF NOT EXISTS event_attendance (
	event_id TEXT NOT NULL,
	user_id TEXT NOT NULL,

	total_seconds INTEGER NOT NULL,

	first_joined_at TEXT,
	last_left_at TEXT,

	was_subscribed INTEGER NOT NULL,

	attended INTEGER NOT NULL,
	no_show INTEGER NOT NULL,
	walk_in INTEGER NOT NULL,

	PRIMARY KEY (event_id, user_id)
);
`
	return nil // Change this to execute the actual schema initialization
}
