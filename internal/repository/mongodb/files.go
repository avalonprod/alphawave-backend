package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FilesRepository struct {
	db *mongo.Collection
}

func NewFilesRepository(db *mongo.Database) *FilesRepository {
	return &FilesRepository{
		db: db.Collection(filesCollection),
	}
}

func (r *FilesRepository) Create(ctx context.Context, input model.File) (string, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.InsertOne(nCtx, input)

	if err != nil {
		return "", err
	}
	ObjectID := res.InsertedID.(primitive.ObjectID)
	id := ObjectID.Hex()
	return id, nil
}

func (r *FilesRepository) RenameFile(ctx context.Context, teamId, fileId, name string, path []string) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ObjectID, err := primitive.ObjectIDFromHex(fileId)
	if err != nil {
		return err
	}

	res, err := r.db.UpdateOne(nCtx, bson.M{"_id": ObjectID, "teamId": teamId}, bson.M{"$set": bson.M{"name": name, "path": path}})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return apperrors.ErrDocumentNotFound
	}

	return nil
}

func (r *FilesRepository) GetFileById(ctx context.Context, teamID, fileID string) (model.File, error) {

	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var file model.File

	ObjectID, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return model.File{}, err
	}

	res := r.db.FindOne(nCtx, bson.M{"_id": ObjectID, "teamId": teamID})
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.File{}, apperrors.ErrDocumentNotFound
		}
		return model.File{}, err
	}

	if err := res.Decode(&file); err != nil {
		return model.File{}, err
	}

	return file, nil

}

func (r *FilesRepository) GetFilesByFolderId(ctx context.Context, teamId, folderId string) ([]model.File, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var files []model.File

	cur, err := r.db.Find(nCtx, bson.M{"teamId": teamId, "folderId": folderId})

	if err != nil {
		return []model.File{}, err
	}

	err = cur.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.File{}, apperrors.ErrDocumentNotFound
		}
		return []model.File{}, err
	}

	if err := cur.All(nCtx, &files); err != nil {
		return []model.File{}, err
	}

	return files, nil
}

func (r *FilesRepository) Delete(ctx context.Context, teamID, fileID string) error {
	nCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()

	ObjectID, err := primitive.ObjectIDFromHex(fileID)

	if err != nil {
		return err
	}

	res, err := r.db.DeleteOne(nCtx, bson.M{"_id": ObjectID, "teamId": teamID})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("no file delete")
	}

	return nil
}
