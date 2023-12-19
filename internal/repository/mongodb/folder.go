package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FolderRepository struct {
	db *mongo.Collection
}

func NewFolderRepository(db *mongo.Database) *FolderRepository {
	return &FolderRepository{
		db: db.Collection(folderCollection),
	}
}

func (r *FolderRepository) CreateFolder(ctx context.Context, teamId string, folder model.Folder) (string, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ObjectId, err := primitive.ObjectIDFromHex(folder.ID)
	if err != nil {
		return "", err
	}

	res, err := r.db.InsertOne(nCtx, bson.M{
		"_id":              ObjectId,
		"teamId":           folder.TeamId,
		"name":             folder.Name,
		"type":             folder.Type,
		"path":             folder.Path,
		"createdAt":        folder.CreatedAt,
		"lastModifiedTime": folder.LastModifiedTime,
		"parentFolder":     folder.ParentFolder,
	})
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *FolderRepository) GetFolderRoot(ctx context.Context, teamId string, folderType string) (model.Folder, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var folder model.Folder
	filter := bson.M{"teamId": teamId, "type": folderType}

	res := r.db.FindOne(nCtx, filter)

	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return folder, apperrors.ErrDocumentNotFound
		}
		return folder, err
	}

	if err := res.Decode(&folder); err != nil {
		return folder, err
	}

	return folder, nil
}

func (r *FolderRepository) GetFolderContentById(ctx context.Context, teamId, id string) ([]model.Folder, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var folders []model.Folder

	cur, err := r.db.Find(nCtx, bson.M{"teamId": teamId, "parentFolder": id})

	if err != nil {
		return []model.Folder{}, err
	}

	err = cur.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.Folder{}, apperrors.ErrDocumentNotFound
		}
		return []model.Folder{}, err
	}

	if err := cur.All(nCtx, &folders); err != nil {
		return []model.Folder{}, err
	}

	return folders, nil
}

func (r *FolderRepository) GetFolderById(ctx context.Context, teamId string, folderId string) (model.Folder, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var folder model.Folder

	ObjectID, err := primitive.ObjectIDFromHex(folderId)
	if err != nil {
		return model.Folder{}, err
	}

	res := r.db.FindOne(nCtx, bson.M{"_id": ObjectID, "teamId": teamId})
	err = res.Err()

	if err != nil {
		return model.Folder{}, err
	}

	if err := res.Decode(&folder); err != nil {
		return model.Folder{}, err
	}

	return folder, nil
}
