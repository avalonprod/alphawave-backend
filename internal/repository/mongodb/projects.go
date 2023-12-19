package mongodb

import "go.mongodb.org/mongo-driver/mongo"

type ProjectsRepository struct {
	db *mongo.Collection
}

func NewProjectsRepository(db *mongo.Database) *ProjectsRepository {
	return &ProjectsRepository{
		db: db.Collection(projectsCollection),
	}
}
