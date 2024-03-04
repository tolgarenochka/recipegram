package dbwizard

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	conn *sqlx.DB
}

func NewConnect() (*Store, error) {
	conn, err := sqlx.ConnectContext(context.Background(), "postgres", "postgresql://localhost:5432/recipegram")
	if err != nil {
		log.Println("Error while db connecting:", err)
		return nil, err
	}
	return &Store{conn: conn}, nil
}

func (s *Store) Quit() error {
	return s.conn.Close()
}
