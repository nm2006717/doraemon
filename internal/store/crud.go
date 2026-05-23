package store

import (
	"database/sql"
	"fmt"
	"strings"
)

func (s *Store) Save(e *Entry) (int64, error) {
	result, err := s.db.Exec(
		`INSERT INTO entries (title, content, category, tags, file_path) VALUES (?, ?, ?, ?, ?)`,
		e.Title, e.Content, e.Category, e.Tags, e.FilePath,
	)
	if err != nil {
		return 0, fmt.Errorf("insert entry: %w", err)
	}
	return result.LastInsertId()
}

func (s *Store) Search(query string) ([]Entry, error) {
	rows, err := s.db.Query(
		`SELECT e.id, e.title, e.content, e.category, e.tags, e.file_path, e.created_at, e.updated_at
		 FROM entries_fts fts
		 JOIN entries e ON e.id = fts.rowid
		 WHERE entries_fts MATCH ?
		 ORDER BY rank`,
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer rows.Close()
	return scanEntries(rows)
}

func (s *Store) List(category string) ([]Entry, error) {
	var rows *sql.Rows
	var err error

	if category == "" {
		rows, err = s.db.Query(
			`SELECT id, title, content, category, tags, file_path, created_at, updated_at
			 FROM entries ORDER BY created_at DESC`,
		)
	} else {
		rows, err = s.db.Query(
			`SELECT id, title, content, category, tags, file_path, created_at, updated_at
			 FROM entries WHERE category = ? ORDER BY created_at DESC`,
			category,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	defer rows.Close()
	return scanEntries(rows)
}

func (s *Store) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM entries WHERE id = ?`, id)
	return err
}

func scanEntries(rows *sql.Rows) ([]Entry, error) {
	var entries []Entry
	for rows.Next() {
		var e Entry
		var filePath sql.NullString
		if err := rows.Scan(&e.ID, &e.Title, &e.Content, &e.Category, &e.Tags, &filePath, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan entry: %w", err)
		}
		e.FilePath = filePath.String
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (s *Store) GetByID(id int64) (*Entry, error) {
	var e Entry
	var filePath sql.NullString
	err := s.db.QueryRow(
		`SELECT id, title, content, category, tags, file_path, created_at, updated_at
		 FROM entries WHERE id = ?`, id,
	).Scan(&e.ID, &e.Title, &e.Content, &e.Category, &e.Tags, &filePath, &e.CreatedAt, &e.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	e.FilePath = filePath.String
	return &e, nil
}

func (s *Store) Update(e *Entry) error {
	_, err := s.db.Exec(
		`UPDATE entries SET title=?, content=?, category=?, tags=?, file_path=?, updated_at=CURRENT_TIMESTAMP
		 WHERE id=?`,
		e.Title, e.Content, e.Category, e.Tags, e.FilePath, e.ID,
	)
	return err
}

func (s *Store) ListPaginated(category string, offset, limit int) ([]Entry, error) {
	var rows *sql.Rows
	var err error

	if category == "" {
		rows, err = s.db.Query(
			`SELECT id, title, content, category, tags, file_path, created_at, updated_at
			 FROM entries ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset,
		)
	} else {
		rows, err = s.db.Query(
			`SELECT id, title, content, category, tags, file_path, created_at, updated_at
			 FROM entries WHERE category = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
			category, limit, offset,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("list paginated: %w", err)
	}
	defer rows.Close()
	return scanEntries(rows)
}

func (s *Store) Count(category string) (int, error) {
	var count int
	var err error
	if category == "" {
		err = s.db.QueryRow(`SELECT COUNT(*) FROM entries`).Scan(&count)
	} else {
		err = s.db.QueryRow(`SELECT COUNT(*) FROM entries WHERE category = ?`, category).Scan(&count)
	}
	return count, err
}

func (s *Store) CountByCategory() (map[string]int, error) {
	rows, err := s.db.Query(`SELECT category, COUNT(*) FROM entries GROUP BY category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var cat string
		var count int
		if err := rows.Scan(&cat, &count); err != nil {
			return nil, err
		}
		result[cat] = count
	}
	return result, rows.Err()
}

func (s *Store) RecentEntries(limit int) ([]Entry, error) {
	rows, err := s.db.Query(
		`SELECT id, title, content, category, tags, file_path, created_at, updated_at
		 FROM entries ORDER BY created_at DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEntries(rows)
}

func BuildFTSQuery(input string) string {
	words := strings.Fields(input)
	for i, w := range words {
		words[i] = w + "*"
	}
	return strings.Join(words, " ")
}
