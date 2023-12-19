package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
)

type PaymentService struct {
	paymentProvider PaymentProvider
	teamsRepository repository.TeamsRepository
}

func NewPaymentService(paymentProvider PaymentProvider, teamsRepository repository.TeamsRepository) *PaymentService {
	return &PaymentService{
		paymentProvider: paymentProvider,
		teamsRepository: teamsRepository,
	}
}

// Create a new customer to be used for payments
func (s *PaymentService) CreateCustomer(name, email, descr string) (*string, error) {
	res, err := s.paymentProvider.CreateCustomer(name, email, fmt.Sprintf("customer with email: %s", email))
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Create new payment method intent
func (s *PaymentService) CreateNewPaymentMethod(ctx context.Context, teamID string) (*string, error) {

	team, err := s.teamsRepository.GetTeamByID(ctx, teamID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return nil, err
		}
		return nil, err
	}

	res, err := s.paymentProvider.NewCard(team.CustomerId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
