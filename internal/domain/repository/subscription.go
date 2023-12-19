package repository

import (
	"context"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, input model.Subscription) error
}
