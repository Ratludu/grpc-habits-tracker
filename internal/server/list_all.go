package server

import (
	"context"
	"encoding/json"

	"github.com/ratludu/grpc-habits-tracker/api"
	"github.com/ratludu/grpc-habits-tracker/internal/habit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) ListHabit(ctx context.Context, request *api.ListHabitRequest) (*api.ListHabitResponse, error) {

	data, err := s.db.GetAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get habits %s", err.Error())
	}

	listHabits := make([]*api.Habit, 0)
	for k, v := range data {
		h := habit.Habit{}
		err = json.Unmarshal([]byte(v), &h)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "coulf not unmarshal habit %s: %s", k, err.Error())
		}
		if h.ID == "" {
			continue
		}
		listHabits = append(listHabits, &api.Habit{
			Id:              string(h.ID),
			Name:            string(h.Name),
			WeeklyFrequency: int32(h.WeeklyFrequency),
		})
	}

	s.lgr.Logf("returned all habit data")

	return &api.ListHabitResponse{
		Habits: listHabits,
	}, nil
}
