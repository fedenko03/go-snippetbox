package models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

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
	fmt.Println(id)
	fmt.Println(int(id))
	return int(id), nil
}
