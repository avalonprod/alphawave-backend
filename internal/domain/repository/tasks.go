package repository

import (
	"context"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type TasksRepository interface {
	CreateTask(ctx context.Context, input model.Task) (string, error)
	UpdatePosition(ctx context.Context, userId string, input []model.Task) error
	GetById(ctx context.Context, userID, taskID string) (model.Task, error)
	GetAll(ctx context.Context, userID string) ([]model.Task, error)
	UpdateById(ctx context.Context, userID string, input model.Task) (model.Task, error)
	ChangeStatus(ctx context.Context, userID, taskID, status string) error
	DeleteAll(ctx context.Context, userID string, status string) error
}
