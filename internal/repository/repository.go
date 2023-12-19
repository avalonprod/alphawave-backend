package repository

import (
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/repository/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	User         repository.UserRepository
	Tasks        repository.TasksRepository
	Roles        repository.RolesRepository
	Projects     repository.ProjectsRepository
	Teams        repository.TeamsRepository
	Members      repository.MemberRepository
	Packages     repository.PackagesRepository
	Files        repository.FilesRepository
	Folder       repository.FolderRepository
	Subscription repository.SubscriptionRepository
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		User:         mongodb.NewUserRepository(db),
		Tasks:        mongodb.NewTasksRepository(db),
		Roles:        mongodb.NewRolesRepository(db),
		Projects:     mongodb.NewProjectsRepository(db),
		Teams:        mongodb.NewTeamsRepository(db),
		Members:      mongodb.NewMemberRepository(db),
		Packages:     mongodb.NewPackagesRepository(db),
		Files:        mongodb.NewFilesRepository(db),
		Folder:       mongodb.NewFolderRepository(db),
		Subscription: mongodb.NewSubscriptionRepository(db),
	}
}
