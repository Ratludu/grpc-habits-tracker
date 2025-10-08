package server

import (
	"context"
	"fmt"

	"github.com/ratludu/grpc-habits-tracker/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) DeleteHabit(ctx context.Context, request *api.DeleteHabitRequest) (*api.DeleteHabitResponse, error) {

	s.lgr.Logf(fmt.Sprintf("Deleting habit Id: %s", request.Id))
	err := s.db.Delete([]byte(request.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "habit could not be deleted %s: %s", request.Id, err.Error())
	}
	s.lgr.Logf("successfully deleted habit data")

	return &api.DeleteHabitResponse{
		Status: "deleted",
	}, nil
}
