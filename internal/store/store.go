package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func New(dataDir string) (*Store, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	dbPath := filepath.Join(dataDir, "doraemon.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			category TEXT NOT NULL DEFAULT 'note',
			tags TEXT NOT NULL DEFAULT '',
			file_path TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE VIRTUAL TABLE IF NOT EXISTS entries_fts USING fts5(
			title, content, tags, content='entries', content_rowid='id'
		);

		CREATE TRIGGER IF NOT EXISTS entries_ai AFTER INSERT ON entries BEGIN
			INSERT INTO entries_fts(rowid, title, content, tags)
			VALUES (new.id, new.title, new.content, new.tags);
		END;

		CREATE TRIGGER IF NOT EXISTS entries_ad AFTER DELETE ON entries BEGIN
			INSERT INTO entries_fts(entries_fts, rowid, title, content, tags)
			VALUES ('delete', old.id, old.title, old.content, old.tags);
		END;

		CREATE TRIGGER IF NOT EXISTS entries_au AFTER UPDATE ON entries BEGIN
			INSERT INTO entries_fts(entries_fts, rowid, title, content, tags)
			VALUES ('delete', old.id, old.title, old.content, old.tags);
			INSERT INTO entries_fts(rowid, title, content, tags)
			VALUES (new.id, new.title, new.content, new.tags);
		END;
	`)
	return err
}
