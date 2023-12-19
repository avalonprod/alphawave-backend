package service

import (
	"context"
	"errors"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
)

type TeamsService struct {
	repository       repository.TeamsRepository
	userRepository   repository.UserRepository
	filesService     FilesServiceI
	memberRepository repository.MemberRepository
	paymentService   PaymentServiceI
	rolesService     RolesService
}

func NewTeamsService(repository repository.TeamsRepository, userRepository repository.UserRepository, memberRepository repository.MemberRepository, paymentService PaymentServiceI, rolesService RolesService, filesService FilesServiceI) *TeamsService {
	return &TeamsService{
		repository:       repository,
		userRepository:   userRepository,
		memberRepository: memberRepository,
		filesService:     filesService,
		rolesService:     rolesService,
		paymentService:   paymentService,
	}
}

func (s *TeamsService) Create(ctx context.Context, userID string, input types.CreateTeamsDTO) (string, error) {
	user, err := s.userRepository.GetUserById(ctx, userID)

	if err != nil {
		return "", err
	}

	customerID, err := s.paymentService.CreateCustomer(input.TeamName, user.Email, input.JobTitle)
	if err != nil {
		return "", err
	}

	team := model.Team{
		TeamName:   input.TeamName,
		JobTitle:   input.JobTitle,
		CustomerId: *customerID,
		OwnerID:    userID,
	}

	id, err := s.repository.CreateTeam(ctx, team)
	if err != nil {
		return "", err
	}

	err = s.rolesService.Create(ctx, id)
	if err != nil {
		return "", err
	}

	if err := s.filesService.CreateRootFolder(ctx, id); err != nil {
		return "", err
	}

	roles := make([]string, 0, 1)

	roles = append(roles, model.ROLE_OWNER)
	if err := s.memberRepository.CreateMember(ctx, id, model.Member{
		TeamID: id,
		UserID: user.ID,
		Email:  user.Email,
		Status: USER_STATUS_ACTIVE,
		Roles:  roles,
	}); err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return "", err
		}
		return "", err
	}

	return id, nil
}

func (s *TeamsService) UpdateTeamSettings(ctx context.Context, teamID string, input types.UpdateTeamSettingsDTO) error {

	if err := s.repository.UpdateTeamSettings(ctx, teamID, model.UpdateTeamSettingsInput{
		LogoURL:               input.LogoURL,
		UserActivityIndicator: input.UserActivityIndicator,
		DisplayLinkPreview:    input.DisplayLinkPreview,
		DisplayFilePreview:    input.DisplayFilePreview,
		EnableGifs:            input.EnableGifs,
		ShowWeekends:          input.ShowWeekends,
		FirstDayOfWeek:        input.FirstDayOfWeek,
	}); err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrTeamNotFound
		}
		return err
	}

	return nil
}

func (s *TeamsService) GetTeamByID(ctx context.Context, userID, teamID string) (model.Team, error) {
	member, err := s.memberRepository.GetMemberByTeamIdAndUserId(ctx, teamID, userID)
	if err != nil {
		return model.Team{}, err
	}

	return s.repository.GetTeamByID(ctx, member.TeamID)
}

func (s *TeamsService) GetTeamByOwnerId(ctx context.Context, ownerId string) (types.TeamsDTO, error) {
	team, err := s.repository.GetTeamByOwnerId(ctx, ownerId)
	if err != nil {
		if errors.Is(err, apperrors.ErrTeamNotFound) {
			return types.TeamsDTO{}, apperrors.ErrTeamNotFound
		}
		return types.TeamsDTO{}, err
	}
	return types.TeamsDTO{
		ID:         team.ID,
		TeamName:   team.TeamName,
		JobTitle:   team.JobTitle,
		OwnerID:    team.OwnerID,
		CustomerId: team.CustomerId,
	}, nil
}

func (s *TeamsService) GetTeamsByUser(ctx context.Context, userID string) ([]model.Team, error) {
	members, err := s.memberRepository.GetMembersByUserID(ctx, userID)

	if err != nil {
		return []model.Team{}, err
	}

	teamsIds := make([]string, 0, len(members))
	for _, member := range members {
		teamsIds = append(teamsIds, member.TeamID)
	}

	teams, err := s.repository.GetTeamsByIds(ctx, teamsIds)
	if err != nil {
		return []model.Team{}, err
	}

	return teams, nil
}

func (s *TeamsService) GetTeamsByIds(ctx context.Context, ids []string) ([]model.Team, error) {
	teams, err := s.repository.GetTeamsByIds(ctx, ids)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return []model.Team{}, apperrors.ErrTeamNotFound
		}
		return []model.Team{}, err
	}
	return teams, nil
}
