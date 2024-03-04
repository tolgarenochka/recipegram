package dbwizard

import (
	"context"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	conn *sqlx.DB
}

func NewConnect() (*Store, error) {
	conn, err := sqlx.ConnectContext(context.Background(), "pgx", "postgresql://localhost:5432/recipegram")
	if err != nil {
		log.Println("Error while db connecting:", err)
		return nil, err
	}
	return &Store{conn: conn}, nil
}

func (s *Store) Quit() error {
	return s.conn.Close()
}
