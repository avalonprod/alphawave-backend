package service

import "github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"

type ProjectsService struct {
	repository repository.ProjectsRepository
}

func NewProjectsService(repository repository.ProjectsRepository) *ProjectsService {
	return &ProjectsService{
		repository: repository,
	}
}
