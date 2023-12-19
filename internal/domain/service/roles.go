package service

import (
	"context"
	"errors"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
)

type RolesService struct {
	repository repository.RolesRepository
}

func NewRolesService(repository repository.RolesRepository) *RolesService {
	return &RolesService{
		repository: repository,
	}
}

func (s *RolesService) Create(ctx context.Context, teamID string) error {
	roles := model.CreateRoles(teamID)

	err := s.repository.Create(ctx, roles)
	return err
}

func (s *RolesService) GetRolesByTeamId(ctx context.Context, teamID string) ([]types.GetRoleDTO, error) {
	roles, err := s.repository.GetRolesByTeamId(ctx, teamID)

	var rolesOutput []types.GetRoleDTO

	for _, item := range roles.Roles {
		rolesOutput = append(rolesOutput, types.GetRoleDTO{
			Role:        item.Role,
			Permissions: item.Permissions,
		})
	}

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return []types.GetRoleDTO{}, apperrors.ErrDocumentNotFound
		}
		return []types.GetRoleDTO{}, err
	}
	return rolesOutput, nil
}

func (s *RolesService) UpdatePermissions(ctx context.Context, teamID string, input []types.UpdatePermissionsDTO) error {
	var inputUpdate []model.UpdatePermissions

	availableRoles := model.GetAvailableRoles()
	rolePermissions := make(map[string]model.Permissions)
	for _, item := range availableRoles {
		rolePermissions[item.Role] = item.Permissions
	}
	for _, item := range input {
		if model.IsAvailableRole(item.Role) {
			rolePermissions[item.Role] = item.Permissions
		}
	}

	for role, permissions := range rolePermissions {
		inputUpdate = append(inputUpdate, model.UpdatePermissions{
			Role:        role,
			Permissions: permissions,
		})
	}

	err := s.repository.UpdatePermissions(ctx, teamID, inputUpdate)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrDocumentNotFound
		}
		return err
	}

	return nil
}
