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

type TeamsRepository struct {
	db *mongo.Collection
}

func NewTeamsRepository(db *mongo.Database) *TeamsRepository {
	return &TeamsRepository{
		db: db.Collection(teamsCollection),
	}
}

func (r *TeamsRepository) CreateTeam(ctx context.Context, input model.Team) (string, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.InsertOne(nCtx, input)

	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), err
}

func (r *TeamsRepository) UpdateTeamSettings(ctx context.Context, teamID string, input model.UpdateTeamSettingsInput) error {
	updateQuery := bson.M{}

	ObjectID, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		return apperrors.ErrInvalidIdFormat
	}
	if input.LogoURL != nil {
		updateQuery["settings.logoUrl"] = *input.LogoURL
	}
	if input.UserActivityIndicator != nil {
		updateQuery["settings.userActivityIndicator"] = *input.UserActivityIndicator
	}
	if input.DisplayLinkPreview != nil {
		updateQuery["settings.displayLinkPreview"] = *input.DisplayLinkPreview
	}
	if input.DisplayFilePreview != nil {
		updateQuery["settings.displayFilePreview"] = *input.DisplayFilePreview
	}
	if input.EnableGifs != nil {
		updateQuery["settings.enableGifs"] = *input.EnableGifs
	}
	if input.ShowWeekends != nil {
		updateQuery["settings.showWeekends"] = *input.ShowWeekends
	}
	if input.FirstDayOfWeek != nil {
		updateQuery["settings.firstDayOfWeek"] = *input.FirstDayOfWeek
	}

	_, err = r.db.UpdateOne(ctx, bson.M{"_id": ObjectID}, bson.M{"$set": updateQuery})
	return err

}

func (r *TeamsRepository) GetTeamsByIds(ctx context.Context, ids []string) ([]model.Team, error) {
	nCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	if len(ids) <= 0 {
		return []model.Team{}, nil
	}

	var teams = make([]model.Team, len(ids))
	var teamsIds = make([]primitive.ObjectID, len(ids))

	for _, id := range ids {

		ObjectID, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			return []model.Team{}, apperrors.ErrInvalidIdFormat
		}
		teamsIds = append(teamsIds, ObjectID)
	}
	filter := bson.M{"_id": bson.M{"$in": teamsIds}}

	cur, err := r.db.Find(nCtx, filter)
	defer cur.Close(nCtx)

	if err != nil {
		return []model.Team{}, err
	}

	err = cur.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.Team{}, apperrors.ErrTeamNotFound
		}
		return []model.Team{}, err
	}

	if err := cur.All(nCtx, &teams); err != nil {
		return []model.Team{}, err
	}
	return teams, nil
}

func (r *TeamsRepository) GetTeamByOwnerId(ctx context.Context, ownerId string) (model.Team, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var team model.Team

	res := r.db.FindOne(nCtx, bson.M{"ownerID": ownerId})

	err := res.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return team, apperrors.ErrTeamNotFound
		}
		return team, err
	}

	if err := res.Decode(&team); err != nil {
		return team, err
	}
	return team, nil
}

func (r *TeamsRepository) GetTeamByID(ctx context.Context, teamID string) (model.Team, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var team model.Team

	ObjectID, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		return model.Team{}, apperrors.ErrInvalidIdFormat
	}

	res := r.db.FindOne(nCtx, bson.M{"_id": ObjectID})

	err = res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return team, apperrors.ErrTeamNotFound
		}

		return team, err
	}

	if err := res.Decode(&team); err != nil {
		return team, err
	}

	return team, nil
}
