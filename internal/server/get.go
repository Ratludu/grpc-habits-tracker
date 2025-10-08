package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ratludu/grpc-habits-tracker/api"
	"github.com/ratludu/grpc-habits-tracker/internal/habit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetHabit(ctx context.Context, request *api.GetHabitRequest) (*api.GetHabitResponse, error) {

	s.lgr.Logf(fmt.Sprintf("Getting habit Id: %s", request.Id))
	data, err := s.db.Get([]byte(request.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "habit not found %s: %s", request.Id, err.Error())
	}

	h := habit.Habit{}
	err = json.Unmarshal(data, &h)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "coulf not unmarshal habit %s: %s", request.Id, err.Error())
	}

	s.lgr.Logf("successfully return habit data")

	return &api.GetHabitResponse{
		Habit: &api.Habit{
			Name:            string(h.Name),
			WeeklyFrequency: int32(h.WeeklyFrequency),
		},
	}, nil
}
