package repository

import (
	"context"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type MemberRepository interface {
	GetMembersByQuery(ctx context.Context, teamID string, query model.GetMembersByQuery) ([]model.Member, error)
	MemberIsDuplicate(ctx context.Context, email string) (bool, error)
	// GetMemberByEmail(ctx context.Context, email string) (model.Member, error)
	UpdateRoles(ctx context.Context, memberId, teamId string, roles []string) error
	GetMemberByToken(ctx context.Context, token string) (model.Member, error)
	GetMemberByTeamIdAndUserId(ctx context.Context, teamID string, userID string) (model.Member, error)
	GetMembersByUserID(ctx context.Context, userID string) ([]model.Member, error)
	CreateMember(ctx context.Context, teamID string, member model.Member) error
	SetStatus(ctx context.Context, memberID string, status string) error
	SetUserID(ctx context.Context, memberID string, userID string) error
	DeleteToken(ctx context.Context, memberID string) error
}
