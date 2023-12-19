package repository

import (
	"context"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type TeamsRepository interface {
	CreateTeam(ctx context.Context, input model.Team) (string, error)
	UpdateTeamSettings(ctx context.Context, teamID string, input model.UpdateTeamSettingsInput) error
	GetTeamByID(ctx context.Context, teamID string) (model.Team, error)
	GetTeamByOwnerId(ctx context.Context, ownerId string) (model.Team, error)
	GetTeamsByIds(ctx context.Context, ids []string) ([]model.Team, error)
}
