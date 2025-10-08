package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/ratludu/grpc-habits-tracker/api"
	"github.com/ratludu/grpc-habits-tracker/internal/database"
	"google.golang.org/grpc"
)

type Server struct {
	lgr Logger
	db  *database.Database
	api.UnimplementedHabitsServer
}

func New(lgr Logger, db *database.Database) *Server {
	return &Server{
		lgr: lgr,
		db:  db,
	}
}

type Logger interface {
	Logf(format string, args ...any)
}

func (s *Server) ListenAndServe(port int) error {
	const addr = "127.0.0.1"

	listener, err := net.Listen("tcp", net.JoinHostPort(addr, strconv.Itoa(port)))
	if err != nil {
		return fmt.Errorf("unable to listen to tcp port %d: %w", port, err)
	}

	grpcServer := grpc.NewServer()
	api.RegisterHabitsServer(grpcServer, s)

	s.lgr.Logf("starting server on port %d\n", port)

	err = grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("error while listening: %w", err)
	}

	return nil
}
