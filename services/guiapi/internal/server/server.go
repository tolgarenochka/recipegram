package server

import (
	"context"
	"github.com/fasthttp/router"
	http "github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	eg  *errgroup.Group
	ctx context.Context

	router *router.Router
	server *HTTPServer
}

func (s *Server) Init() {
	s.router = s.initRouter()
	s.server.serverHTTP.Handler = s.router.Handler
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
		return s.server.serverHTTP.ListenAndServe("localhost:8080")
	})

	return s.eg.Wait()
}
