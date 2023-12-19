package repository

import (
	"context"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type RolesRepository interface {
	Create(ctx context.Context, input model.TeamRoles) error
	GetRolesByTeamId(ctx context.Context, teamID string) (model.TeamRoles, error)
	UpdatePermissions(ctx context.Context, teamID string, input []model.UpdatePermissions) error
}
