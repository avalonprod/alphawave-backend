package service

import (
	"context"
	"errors"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
)

type TasksService struct {
	repository repository.TasksRepository
}

func NewTasksService(repository repository.TasksRepository) *TasksService {
	return &TasksService{
		repository: repository,
	}
}

const (
	STATUS_DELITE string = "del"
	STATUS_DONE   string = "done"
	STATUS_ACTIVE string = "active"
)

func (s *TasksService) Create(ctx context.Context, userID string, input types.TasksCreateDTO) (string, error) {

	task := model.Task{
		UserID:   userID,
		Title:    input.Title,
		Status:   STATUS_ACTIVE,
		Priority: input.Priority,
		Index:    input.Index,
	}
	id, err := s.repository.CreateTask(ctx, task)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *TasksService) UpdatePosition(ctx context.Context, userId string, input []types.UpdatePositionDTO) error {
	var tasks []model.Task

	for i, item := range input {
		tasks = append(tasks, model.Task{
			ID:    item.Id,
			Index: i,
		})
	}
	err := s.repository.UpdatePosition(ctx, userId, tasks)
	if err != nil {
		return err
	}
	return nil
}

func (s *TasksService) GetById(ctx context.Context, userID, taskID string) (types.TaskDTO, error) {
	task, err := s.repository.GetById(ctx, userID, taskID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return types.TaskDTO{}, apperrors.ErrDocumentNotFound
		}
		return types.TaskDTO{}, err
	}

	return types.TaskDTO{
		ID:       task.ID,
		Title:    task.Title,
		Priority: task.Priority,
		Index:    task.Index,
	}, nil
}

func (s *TasksService) GetAll(ctx context.Context, userID string) (types.TasksDTO, error) {
	tasksIn, err := s.repository.GetAll(ctx, userID)

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return types.TasksDTO{}, err
		}
		return types.TasksDTO{}, err
	}

	tasks := types.TasksDTO{
		ActiveTasks:   []types.TaskDTO{},
		FinishedTasks: []types.TaskDTO{},
		DeletedTasks:  []types.TaskDTO{},
	}

	for i := range tasksIn {
		task := types.TaskDTO{
			ID:       tasksIn[i].ID,
			Title:    tasksIn[i].Title,
			Priority: tasksIn[i].Priority,
			Index:    tasksIn[i].Index,
		}
		if tasksIn[i].Status == STATUS_ACTIVE {
			tasks.ActiveTasks = append(tasks.ActiveTasks, task)
		}
		if tasksIn[i].Status == STATUS_DONE {
			tasks.FinishedTasks = append(tasks.FinishedTasks, task)
		}
		if tasksIn[i].Status == STATUS_DELITE {
			tasks.DeletedTasks = append(tasks.DeletedTasks, task)
		}

	}

	return tasks, nil
}

func (s *TasksService) UpdateById(ctx context.Context, userID string, input types.UpdateTaskDTO) (types.TaskDTO, error) {

	task, err := s.repository.UpdateById(ctx, userID, model.Task{
		ID:       input.ID,
		Title:    input.Title,
		Priority: input.Priority,
		Index:    input.Index,
	})

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return types.TaskDTO{}, apperrors.ErrDocumentNotFound
		}
		return types.TaskDTO{}, err
	}
	return types.TaskDTO{
		ID:       task.ID,
		Title:    task.Title,
		Priority: task.Priority,
		Index:    task.Index,
	}, err
}

func (s *TasksService) DeleteTaskById(ctx context.Context, userID, taskID string) error {
	err := s.repository.ChangeStatus(ctx, userID, taskID, STATUS_DELITE)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrDocumentNotFound
		}
		return err
	}
	return nil
}

func (s *TasksService) FinishedTaskById(ctx context.Context, userID, taskID string) error {
	err := s.repository.ChangeStatus(ctx, userID, taskID, STATUS_DONE)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrDocumentNotFound
		}
		return err
	}
	return nil
}

func (s *TasksService) DeleteAll(ctx context.Context, userID string) error {
	err := s.repository.DeleteAll(ctx, userID, STATUS_DELITE)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrDocumentNotFound
		}
		return err
	}
	return nil
}

func (s *TasksService) ClearAll(ctx context.Context, userID string) error {
	err := s.repository.DeleteAll(ctx, userID, STATUS_DONE)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrDocumentNotFound
		}
		return err
	}
	return nil
}

func (s *TasksService) UndoTask(ctx context.Context, userID string, taskID string) error {
	err := s.repository.ChangeStatus(ctx, userID, taskID, STATUS_ACTIVE)

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrDocumentNotFound
		}
		return err
	}
	return nil
}
