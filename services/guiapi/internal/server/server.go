package server

import (
	"context"
	"github.com/fasthttp/router"
	"github.com/jmoiron/sqlx"
	http "github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	eg  *errgroup.Group
	ctx context.Context

	router   *router.Router
	server   *HTTPServer
	dbWizard dbWiz
}

type dbWiz struct {
	dbWizard *sqlx.DB
}

func (s *Server) Init(dbW *sqlx.DB) {
	s.router = s.initRouter()
	s.server.serverHTTP.Handler = s.router.Handler
	s.dbWizard = dbWiz{
		dbWizard: dbW,
	}
}

type HTTPServer struct {
	serverHTTP *http.Server
}

func NewServer() *Server {
	return &Server{
		server: &HTTPServer{
			serverHTTP: &http.Server{
				Name: "http server",
			},
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.eg, s.ctx = errgroup.WithContext(ctx)

	s.eg.Go(func() error {
		return s.server.serverHTTP.ListenAndServe("0.0.0.0:8080")
	})

	return s.eg.Wait()
}
