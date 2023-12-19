package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
)

type SubscriptionService struct {
	UserService     UserServiceI
	teamsRepository repository.TeamsRepository
	packagesService PackagesServiceI
	paymentProvider PaymentProvider
	repository      repository.SubscriptionRepository
}

func NewSubscriptionService(userService UserServiceI, teamsRepository repository.TeamsRepository, packagesService PackagesServiceI, repository repository.SubscriptionRepository, paymantProvider PaymentProvider) *SubscriptionService {
	return &SubscriptionService{
		UserService:     userService,
		teamsRepository: teamsRepository,
		repository:      repository,
		paymentProvider: paymantProvider,
		packagesService: packagesService,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, userID string, packageID string, teamID string) error {
	user, err := s.UserService.GetUserById(ctx, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrUserNotFound
		}
		return err
	}

	team, err := s.teamsRepository.GetTeamByID(ctx, teamID)
	if err != nil {
		return err
	}

	pkg, err := s.packagesService.GetById(ctx, packageID)

	if err != nil {
		return err
	}
	fmt.Printf("customerID: %s", team.CustomerId)

	subId, err := s.paymentProvider.CreateSubscription(team.CustomerId, pkg.StripePriceId)

	if err != nil {
		return err
	}

	err = s.repository.CreateSubscription(ctx, model.Subscription{
		StripeSubId: *subId,
		TeamID:      team.ID,
		UserInfo: model.UserInfoShort{
			ID:        userID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		},
		Status: model.OrderStatusCreated,
	})

	if err != nil {
		return err
	}

	return nil
}
