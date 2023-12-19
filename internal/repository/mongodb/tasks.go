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

type TasksRepository struct {
	db *mongo.Collection
}

func NewTasksRepository(db *mongo.Database) *TasksRepository {
	return &TasksRepository{
		db: db.Collection(tasksCollection),
	}
}

func (r *TasksRepository) CreateTask(ctx context.Context, input model.Task) (string, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	res, err := r.db.InsertOne(nCtx, input)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *TasksRepository) UpdatePosition(ctx context.Context, userId string, input []model.Task) error {
	nCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	var bulkWriteOptions []mongo.WriteModel

	for _, item := range input {
		ObjectID, err := primitive.ObjectIDFromHex(item.ID)
		if err != nil {
			return err
		}

		updateDocument := mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": ObjectID}).SetUpdate(bson.M{"$set": bson.M{"index": item.Index}})
		bulkWriteOptions = append(bulkWriteOptions, updateDocument)
	}

	_, err := r.db.BulkWrite(nCtx, bulkWriteOptions)

	if err != nil {
		return err
	}
	return nil
}

func (r *TasksRepository) GetById(ctx context.Context, userID, taskID string) (model.Task, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var task model.Task

	ObjectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return model.Task{}, err
	}

	filter := bson.M{"_id": ObjectID, "userID": userID}

	res := r.db.FindOne(nCtx, filter)

	err = res.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return task, apperrors.ErrDocumentNotFound
		}
		return task, err
	}

	if err := res.Decode(&task); err != nil {
		return task, err
	}

	return task, nil
}

func (r *TasksRepository) GetAll(ctx context.Context, userID string) ([]model.Task, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var tasks []model.Task

	cur, err := r.db.Find(nCtx, bson.M{"userID": userID})

	if err != nil {
		return []model.Task{}, err
	}

	err = cur.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.Task{}, apperrors.ErrDocumentNotFound
		}
		return []model.Task{}, err
	}

	if err := cur.All(nCtx, &tasks); err != nil {
		return []model.Task{}, err
	}

	return tasks, nil
}

func (r *TasksRepository) UpdateById(ctx context.Context, userID string, input model.Task) (model.Task, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ObjectID, err := primitive.ObjectIDFromHex(input.ID)
	if err != nil {
		return model.Task{}, err
	}

	filter := bson.M{"_id": ObjectID, "userID": userID}

	inputByte, err := bson.Marshal(input)

	if err != nil {
		return model.Task{}, fmt.Errorf("failed to marshal document. error: %s", err)
	}

	var updateObj bson.M

	err = bson.Unmarshal(inputByte, &updateObj)
	if err != nil {
		return model.Task{}, fmt.Errorf("failed to unmarshal document. error: %s", err)
	}

	delete(updateObj, "_id")

	update := bson.M{
		"$set": updateObj,
	}

	result, err := r.db.UpdateOne(nCtx, filter, update)

	if err != nil {
		return model.Task{}, fmt.Errorf("failed to execute query. error: %s", err)
	}
	if result.MatchedCount == 0 {
		return model.Task{}, apperrors.ErrDocumentNotFound
	}

	return input, nil
}

func (r *TasksRepository) ChangeStatus(ctx context.Context, userID, taskID, status string) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ObjectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": ObjectID, "userID": userID}

	res, err := r.db.UpdateOne(nCtx, filter, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return err
	}
	if res.ModifiedCount <= 0 {
		return apperrors.ErrDocumentNotFound
	}
	return nil
}

func (r *TasksRepository) DeleteAll(ctx context.Context, userID string, status string) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"userID": userID, "status": status}
	res, err := r.db.DeleteMany(nCtx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount <= 0 {
		return apperrors.ErrDocumentNotFound
	}

	return nil
}
