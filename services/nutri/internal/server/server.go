package server

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	eg       *errgroup.Group
	ctx      context.Context
	dbWizard dbWiz
	consumer Cons
}

type dbWiz struct {
	dbWizard *sqlx.DB
}

func Init(dbW *sqlx.DB, consumer sarama.Consumer) *Server {
	s := Server{}
	s.dbWizard = dbWiz{
		dbWizard: dbW,
	}
	s.consumer.Consumer = consumer

	return &s
}

func (s *Server) Run(ctx context.Context) error {
	s.eg, s.ctx = errgroup.WithContext(ctx)

	s.eg.Go(func() error {
		return s.consumer.Run(s.ctx)
	})

	return s.eg.Wait()
}
