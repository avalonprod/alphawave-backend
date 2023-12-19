package mongodb

import (
	"context"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscriptionRepository struct {
	db *mongo.Collection
}

func NewSubscriptionRepository(db *mongo.Database) *SubscriptionRepository {
	return &SubscriptionRepository{
		db: db.Collection(subscriptionCollection),
	}
}

func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, input model.Subscription) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err := r.db.InsertOne(nCtx, input); err != nil {
		return err
	}
	return nil
}
