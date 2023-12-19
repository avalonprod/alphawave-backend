package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/codegenerator"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/tokengenerator"
)

const (
	USER_STATUS_ACTIVE   = "ACTIVE"
	USER_STATUS_INACTIVE = "INACTIVE"
	USER_STATUS_PENDING  = "PENDING"
)

type MemberService struct {
	repository      repository.MemberRepository
	userRepository  repository.UserRepository
	userService     UserServiceI
	teamsRepository repository.TeamsRepository
	teamsService    TeamsServiceI
	emailService    *EmailService
	codeGenerator   *codegenerator.CodeGenerator
	tokenGenerator  *tokengenerator.TokenGenerator
	apiUrl          string
}

func NewMemberService(repository repository.MemberRepository, userRepository repository.UserRepository, teamsRepository repository.TeamsRepository, codeGenerator *codegenerator.CodeGenerator, tokenGenerator *tokengenerator.TokenGenerator, teamsService TeamsServiceI, emailService *EmailService, userService UserServiceI, apiUrl string) *MemberService {
	return &MemberService{
		repository:      repository,
		userRepository:  userRepository,
		userService:     userService,
		codeGenerator:   codeGenerator,
		tokenGenerator:  tokenGenerator,
		teamsRepository: teamsRepository,
		teamsService:    teamsService,
		emailService:    emailService,
		apiUrl:          apiUrl,
	}
}

func (s *MemberService) MemberSignUp(ctx context.Context, token string, input types.MemberSignUpDTO) error {
	member, err := s.repository.GetMemberByToken(ctx, token)

	if err != nil {
		if errors.Is(err, apperrors.ErrMemberNotFound) {
			return apperrors.ErrMemberNotFound
		}
		return err
	}

	err = s.userService.SignUp(ctx, types.UserSignUpDTO{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		JobTitle:  input.JobTitle,
		Email:     member.Email,
		Password:  input.Password,
	})

	if err != nil {
		return err
	}

	if err := s.SetStatus(ctx, member.ID, USER_STATUS_ACTIVE); err != nil {
		return err
	}

	user, err := s.userRepository.GetUserByEmail(ctx, member.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrUserNotFound
		}
		return err
	}

	if err := s.SetUserID(ctx, member.ID, user.ID); err != nil {
		return err
	}

	if err := s.repository.DeleteToken(ctx, member.ID); err != nil {
		return err
	}

	return nil
}

func (s *MemberService) GetMembersByQuery(ctx context.Context, teamID string, query types.GetMembersByQuery) ([]types.MemberDTO, error) {
	members, err := s.repository.GetMembersByQuery(ctx, teamID, model.GetMembersByQuery{PaginationQuery: model.PaginationQuery(query.PaginationQuery)})

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return []types.MemberDTO{}, apperrors.ErrDocumentNotFound
		}
		return []types.MemberDTO{}, err
	}
	var userIds = []string{}

	for _, member := range members {
		if member.Status == USER_STATUS_PENDING {
			continue
		}
		userIds = append(userIds, member.UserID)
	}

	users, err := s.userRepository.GetUsersByIds(ctx, userIds)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return []types.MemberDTO{}, apperrors.ErrDocumentNotFound
		}
		return []types.MemberDTO{}, err
	}

	var membersOutput = []types.MemberDTO{}

	for i := range members {
		var j int = 0
		for j = range users {
			if len(userIds) > 0 && users[j].ID == members[i].UserID {
				membersOutput = append(membersOutput, types.MemberDTO{
					MemberID:  members[i].ID,
					FirstName: users[i].FirstName,
					LastName:  users[i].LastName,
					Email:     users[i].Email,
					Roles:     members[i].Roles,
					Status:    members[i].Status,
				})
			}
		}

		if members[i].Status == USER_STATUS_PENDING || members[i].Status == USER_STATUS_INACTIVE {
			membersOutput = append(membersOutput, types.MemberDTO{
				MemberID:  members[i].ID,
				FirstName: "",
				LastName:  "",
				Email:     members[i].Email,
				Roles:     members[i].Roles,
				Status:    members[i].Status,
			})
		}

	}

	return membersOutput, nil
}

func (s *MemberService) GetMemberByTeamIdAndUserId(ctx context.Context, teamID string, userID string) (model.Member, error) {
	member, err := s.repository.GetMemberByTeamIdAndUserId(ctx, teamID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrMemberNotFound) {
			return model.Member{}, apperrors.ErrMemberNotFound
		}
		return model.Member{}, err
	}
	return member, nil
}

func (s *MemberService) UserInvite(ctx context.Context, teamID string, email string, role string) error {

	isDuplicate, err := s.repository.MemberIsDuplicate(ctx, email)
	if err != nil {
		return err
	}

	if isDuplicate {
		return apperrors.ErrUserAlreadyExists
	}

	token := s.codeGenerator.RandomSecret(20)

	roles := make([]string, 0, 1)
	if model.IsAvailableRole(role) {
		roles = append(roles, role)
	} else {
		return errors.New("this role not available")
	}

	team, err := s.teamsRepository.GetTeamByID(ctx, teamID)

	if err != nil {
		return err
	}

	err = s.repository.CreateMember(ctx, teamID, model.Member{
		TeamID:      teamID,
		Email:       email,
		VerifyToken: token,
		Status:      USER_STATUS_PENDING,
		Roles:       roles,
	})

	if err != nil {
		return err
	}

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrDocumentNotFound
		}
		return err
	}

	invationURL := fmt.Sprintf("%s/api/v1/members/accept_invite/%s", s.apiUrl, token)

	err = s.emailService.SendUserInvite(UserInviteInput{
		Email:    email,
		TeamName: team.TeamName,
		URL:      invationURL,
	})

	return err
}

func (s *MemberService) AcceptInvite(ctx context.Context, token string) (string, error) {
	member, err := s.repository.GetMemberByToken(ctx, token)

	if err != nil {
		if errors.Is(err, apperrors.ErrMemberNotFound) {
			return "", apperrors.ErrMemberNotFound
		}
		return "", err
	}
	user, err := s.userRepository.GetUserByEmail(ctx, member.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return member.VerifyToken, apperrors.ErrUserNotFound
		}
		return "", err
	}
	if err := s.SetStatus(ctx, member.ID, USER_STATUS_ACTIVE); err != nil {
		return "", err
	}
	if err := s.SetUserID(ctx, member.ID, user.ID); err != nil {
		return "", err
	}

	if err := s.repository.DeleteToken(ctx, member.ID); err != nil {
		return "", err
	}

	return "", nil
}

func (s *MemberService) SetStatus(ctx context.Context, memberID string, status string) error {
	err := s.repository.SetStatus(ctx, memberID, status)
	if err != nil {
		if errors.Is(err, apperrors.ErrMemberNotFound) {
			return apperrors.ErrMemberNotFound
		}
		return err
	}

	return nil
}

func (s *MemberService) SetUserID(ctx context.Context, memberID string, userID string) error {
	err := s.repository.SetUserID(ctx, memberID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrMemberNotFound) {
			return apperrors.ErrMemberNotFound
		}
		return err
	}
	return nil
}
func (s *MemberService) UpdateRoles(ctx context.Context, memberId, teamId string, roles []string) error {
	for _, role := range roles {
		res := model.IsAvailableRole(role)
		if res == false {
			return apperrors.ErrRoleIsNotAvailable
		}
	}
	err := s.repository.UpdateRoles(ctx, memberId, teamId, roles)
	if err != nil {
		return err
	}
	return nil
}

// func (s *MemberService) ChangeUserEmail(ctx context.Context, memberID, teamID, email string) error {

// 	member := s.repository.GetMemberByTeamIdAndUserId()
// }
