package server

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ratludu/grpc-habits-tracker/api"
	"github.com/ratludu/grpc-habits-tracker/internal/habit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateHabit(ctx context.Context, request *api.CreateHabitRequest) (*api.CreateHabitResponse, error) {
	s.lgr.Logf("CreateHabit request received: %s", request)

	var freq uint
	if request.WeeklyFrequency != nil {
		freq = uint(*request.WeeklyFrequency)
	}

	h := habit.Habit{
		Name:            habit.Name(request.Name),
		WeeklyFrequency: habit.WeeklyFrequency(freq),
	}

	createdHabit, err := habit.Create(ctx, h)
	if err != nil {
		var invalidErr habit.InvalidInputError
		if errors.As(err, &invalidErr) {
			return nil, status.Errorf(codes.Internal, "cannot save habit %v: %s", h, err.Error())
		}
	}

	jHabit, err := json.Marshal(createdHabit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot save habit in database %v: %s", h, err.Error())
	}

	err = s.db.Set([]byte(createdHabit.ID), jHabit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot save habit in database %v: %s", h, err.Error())
	}

	s.lgr.Logf("Habit %s sucessfully registered", createdHabit.ID)

	return &api.CreateHabitResponse{
		Habit: &api.Habit{
			Id:              string(createdHabit.ID),
			Name:            string(createdHabit.Name),
			WeeklyFrequency: int32(createdHabit.WeeklyFrequency),
		},
	}, nil
}
