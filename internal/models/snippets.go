package models

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *pgx.Conn
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	row := m.DB.QueryRow(context.Background(),
		"INSERT INTO snippets  (title, content, created, expires) VALUES ($1,$2, now(), now() + INTERVAL '1 day' * $3) RETURNING id;",
		title, content, expires)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	s := &Snippet{}
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > now() AND id = $1`
	err := m.DB.QueryRow(context.Background(), stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
WHERE expires > NOW() ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	snippets := []*Snippet{}
	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
