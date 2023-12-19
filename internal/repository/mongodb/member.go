package mongodb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MemberRepository struct {
	db *mongo.Collection
}

func NewMemberRepository(db *mongo.Database) *MemberRepository {
	return &MemberRepository{
		db: db.Collection(memberCollection),
	}
}

func (r *MemberRepository) CreateMember(ctx context.Context, teamID string, member model.Member) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	fmt.Print(member)
	if _, err := r.db.InsertOne(nCtx, member); err != nil {
		return err
	}
	return nil
}

// func (r *MemberRepository) GetMembers(ctx context.Context, teamID string) ([]model.Member, error) {
// 	nCtx, cancel := context.WithTimeout(ctx, 50*time.Second)
// 	defer cancel()

// 	var members []model.Member

// 	cur, err := r.db.Find(nCtx, bson.M{"teamID": teamID})

// 	if err != nil {
// 		return []model.Member{}, err
// 	}

// 	err = cur.Err()

// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return []model.Member{}, apperrors.ErrDocumentNotFound
// 		}
// 		return []model.Member{}, err
// 	}

// 	if err := cur.All(nCtx, &members); err != nil {
// 		return []model.Member{}, err
// 	}

// 	return members, nil
// }

func (r *MemberRepository) GetMembersByQuery(ctx context.Context, teamID string, query model.GetMembersByQuery) ([]model.Member, error) {
	nCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()
	paginationOpts := getPaginationOptions(&query.PaginationQuery)

	var members []model.Member

	cur, err := r.db.Find(nCtx, bson.M{"teamID": teamID}, paginationOpts)
	defer cur.Close(nCtx)

	if err != nil {
		return []model.Member{}, err
	}

	err = cur.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.Member{}, apperrors.ErrDocumentNotFound
		}
		return []model.Member{}, err
	}

	if err := cur.All(nCtx, &members); err != nil {
		return []model.Member{}, err
	}
	return members, nil
}

func (r *MemberRepository) GetMemberByTeamIdAndUserId(ctx context.Context, teamID string, userID string) (model.Member, error) {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var member model.Member

	filter := bson.M{"teamID": teamID, "userID": userID}

	res := r.db.FindOne(nCtx, filter)

	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return member, apperrors.ErrMemberNotFound
		}

		return member, err

	}
	if err := res.Decode(&member); err != nil {
		return member, err
	}

	return member, nil
}

func (r *MemberRepository) GetMembersByUserID(ctx context.Context, userID string) ([]model.Member, error) {
	nCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var members []model.Member

	cur, err := r.db.Find(nCtx, bson.M{"userID": userID})

	if err != nil {
		return []model.Member{}, err
	}

	err = cur.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.Member{}, apperrors.ErrDocumentNotFound
		}
		return []model.Member{}, err
	}

	if err := cur.All(nCtx, &members); err != nil {
		return []model.Member{}, err
	}

	return members, nil
}

func (r *MemberRepository) MemberIsDuplicate(ctx context.Context, email string) (bool, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"email": email}

	count, err := r.db.CountDocuments(nCtx, filter)

	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (r *MemberRepository) GetMemberByToken(ctx context.Context, token string) (model.Member, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var member model.Member

	filter := bson.M{"verifyToken": token}

	res := r.db.FindOne(nCtx, filter)

	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return member, apperrors.ErrMemberNotFound
		}

		return member, err

	}
	if err := res.Decode(&member); err != nil {
		return member, err
	}
	return member, nil
}

func (r *MemberRepository) GetMemberById(ctx context.Context, memberID string, teamID string) (model.Member, error) {
	// todo
	return model.Member{}, nil
}

func (r *MemberRepository) SetStatus(ctx context.Context, memberID string, status string) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ObjectID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return err
	}
	_, err = r.db.UpdateOne(nCtx, bson.M{"_id": ObjectID}, bson.M{"$set": bson.M{"status": status}})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperrors.ErrMemberNotFound
		}
		return err
	}
	return nil
}

func (r *MemberRepository) SetUserID(ctx context.Context, memberID string, userID string) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ObjectID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return err
	}
	_, err = r.db.UpdateOne(nCtx, bson.M{"_id": ObjectID}, bson.M{"$set": bson.M{"userID": userID}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperrors.ErrMemberNotFound
		}
		return err
	}
	return nil
}

func (r *MemberRepository) DeleteToken(ctx context.Context, memberID string) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ObjectID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return err
	}
	_, err = r.db.UpdateOne(nCtx, bson.M{"_id": ObjectID}, bson.M{"$set": bson.M{"verifyToken": ""}})

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperrors.ErrMemberNotFound
		}
		return err
	}

	return nil
}

func (r *MemberRepository) UpdateRoles(ctx context.Context, memberId, teamId string, roles []string) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	ObjectID, err := primitive.ObjectIDFromHex(memberId)
	if err != nil {
		return err
	}
	_, err = r.db.UpdateOne(nCtx, bson.M{"_id": ObjectID, "teamId": teamId}, bson.M{"$set": bson.M{"roles": roles}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperrors.ErrMemberNotFound
		}
		return err
	}
	return nil
}
