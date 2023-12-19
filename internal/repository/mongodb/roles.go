package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RolesRepository struct {
	db *mongo.Collection
}

func NewRolesRepository(db *mongo.Database) *RolesRepository {
	return &RolesRepository{
		db: db.Collection(rolesCollection),
	}
}

func (r *RolesRepository) Create(ctx context.Context, input model.TeamRoles) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.db.InsertOne(nCtx, input)
	if err != nil {
		return err
	}
	return nil
}

func (r *RolesRepository) GetRolesByTeamId(ctx context.Context, teamID string) (model.TeamRoles, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var roles model.TeamRoles

	res := r.db.FindOne(nCtx, bson.M{"teamID": teamID})

	err := res.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.TeamRoles{}, apperrors.ErrDocumentNotFound
		}
		return model.TeamRoles{}, err
	}

	if err := res.Decode(&roles); err != nil {
		return model.TeamRoles{}, err
	}

	return roles, err
}

func (r *RolesRepository) UpdatePermissions(ctx context.Context, teamID string, input []model.UpdatePermissions) error {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(nCtx, bson.M{"teamID": teamID}, bson.M{"$set": bson.M{"roles": input}})

	if err != nil {

		return err
	}
	return nil
}
