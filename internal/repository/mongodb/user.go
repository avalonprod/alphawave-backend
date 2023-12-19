package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		db: db.Collection(usersCollection),
	}
}

func (r *UserRepository) Create(ctx context.Context, input model.User) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err := r.db.InsertOne(nCtx, input); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByСredentials(ctx context.Context, email, password string) (model.User, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user model.User
	filter := bson.M{"email": email, "password": password}

	res := r.db.FindOne(nCtx, filter)

	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, apperrors.ErrUserNotFound
		}
		return user, err
	}

	if err := res.Decode(&user); err != nil {
		return user, err
	}

	return user, err
}

func (r *UserRepository) DeleteUserByEmail(ctx context.Context, email string) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.db.DeleteOne(nCtx, bson.M{"email": email})

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return apperrors.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user model.User
	filter := bson.M{"email": email}

	res := r.db.FindOne(nCtx, filter)

	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, apperrors.ErrUserNotFound
		}
		return user, err
	}

	if err := res.Decode(&user); err != nil {
		return user, err
	}

	return user, err
}

func (r *UserRepository) GetUserById(ctx context.Context, userID string) (model.User, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user model.User
	ObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return model.User{}, err
	}
	filter := bson.M{"_id": ObjectID}

	res := r.db.FindOne(nCtx, filter)

	err = res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, apperrors.ErrUserNotFound
		}

		return user, err

	}
	if err := res.Decode(&user); err != nil {
		return user, err
	}

	return user, err
}

func (r *UserRepository) GetUsersByQuery(ctx context.Context, ids []string, query model.GetUsersByQuery) ([]model.User, error) {
	nCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	paginationOpts := getPaginationOptions(&query.PaginationQuery)

	var users = make([]model.User, len(ids))
	var userIds = make([]primitive.ObjectID, len(ids))

	for _, id := range ids {
		ObjectID, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			return []model.User{}, err
		}
		userIds = append(userIds, ObjectID)
	}

	filter := bson.M{"_id": bson.M{"$in": userIds}}

	cur, err := r.db.Find(nCtx, filter, paginationOpts)
	defer cur.Close(nCtx)

	if err != nil {
		return []model.User{}, err
	}

	err = cur.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.User{}, apperrors.ErrDocumentNotFound
		}
		return []model.User{}, err
	}

	if err := cur.Decode(users); err != nil {
		return []model.User{}, err
	}
	return users, nil
}

func (r *UserRepository) GetUsersByIds(ctx context.Context, ids []string) ([]model.User, error) {
	nCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	if len(ids) <= 0 {
		return []model.User{}, nil
	}

	var users = make([]model.User, len(ids))
	var userIds = make([]primitive.ObjectID, len(ids))

	for _, id := range ids {

		ObjectID, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			return []model.User{}, err
		}
		userIds = append(userIds, ObjectID)
	}
	filter := bson.M{"_id": bson.M{"$in": userIds}}

	cur, err := r.db.Find(nCtx, filter)
	defer cur.Close(nCtx)

	if err != nil {
		return []model.User{}, err
	}

	err = cur.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.User{}, apperrors.ErrDocumentNotFound
		}
		return []model.User{}, err
	}

	if err := cur.All(nCtx, &users); err != nil {
		return []model.User{}, err
	}
	return users, nil
}

func (r *UserRepository) ChangeVerificationCode(ctx context.Context, email string, input model.UserVerificationPayload) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := r.db.UpdateOne(nCtx, bson.M{"email": email}, bson.M{"$set": bson.M{"verification.verificationCode": input.VerificationCode, "verification.verificationCodeExpiresTime": input.VerificationCodeExpiresTime}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperrors.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *UserRepository) SetSession(ctx context.Context, userID string, session model.Session, lastVisitTime time.Time) error {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	ObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	_, err = r.db.UpdateOne(nCtx, bson.M{"_id": ObjectID}, bson.M{"$set": bson.M{"session": session, "lastVisitTime": lastVisitTime}})

	return err
}

func (r *UserRepository) Verify(ctx context.Context, verificationCode string) error {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(nCtx, bson.M{"verification.verificationCode": verificationCode}, bson.M{"$set": bson.M{"verification.verified": true, "verification.verificationCode": ""}})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) ChangePassword(ctx context.Context, userID, newPassword, oldPassword string) error {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	ObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	res, err := r.db.UpdateOne(nCtx, bson.M{"_id": ObjectID, "password": oldPassword}, bson.M{"$set": bson.M{"password": newPassword}})

	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdateUserInfo(ctx context.Context, userID string, input model.UpdateUserInfoInput) error {
	updateQuery := bson.M{}

	if input.FirstName != nil {
		updateQuery["firstName"] = *input.FirstName
	}
	if input.LastName != nil {
		updateQuery["lastName"] = input.LastName
	}
	if input.JobTitle != nil {
		updateQuery["jobTitle"] = input.JobTitle
	}
	if input.Email != nil {
		updateQuery["email"] = input.Email
	}
	ObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = r.db.UpdateOne(ctx, bson.M{"_id": ObjectID}, bson.M{"$set": updateQuery})
	return err
}

func (r *UserRepository) UpdateUserSettings(ctx context.Context, userID string, input model.UpdateUserSettingsInput) error {
	updateQuery := bson.M{}

	if input.UserIcon != nil {
		updateQuery["settings.userIcon"] = input.UserIcon
	}
	if input.BannerImage != nil {
		updateQuery["settings.bannerImage"] = input.BannerImage
	}
	if input.TimeZone != nil {
		updateQuery["settings.timeZone"] = input.TimeZone
	}
	if input.DateFormat != nil {
		updateQuery["settings.dateFormat"] = input.DateFormat
	}
	if input.TimeFormat != nil {
		updateQuery["settings.timeFormat"] = input.TimeFormat
	}
	ObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = r.db.UpdateOne(ctx, bson.M{"_id": ObjectID}, bson.M{"$set": updateQuery})
	return err
}

func (r *UserRepository) SetForgotPassword(ctx context.Context, email string, input model.ForgotPasswordPayload) error {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := r.db.UpdateOne(nCtx, bson.M{"email": email}, bson.M{"$set": bson.M{"forgotPasswordToken": input}})

	return err
}

func (r *UserRepository) GetUserByVerificationCode(ctx context.Context, hash string) (model.User, error) {
	var user model.User
	filter := bson.M{"verification.verificationCode": hash}

	res := r.db.FindOne(ctx, filter)

	err := res.Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, apperrors.ErrUserNotFound
		}

		return user, err

	}
	if err := res.Decode(&user); err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (model.User, error) {
	nCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var user model.User

	if err := r.db.FindOne(nCtx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresTime":  bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, apperrors.ErrUserNotFound
		}

		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepository) RemoveSession(ctx context.Context, userID string) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	ObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	res, err := r.db.UpdateOne(nCtx, bson.M{"_id": ObjectID}, bson.M{"$unset": bson.M{"session": ""}})

	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return apperrors.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) IsDuplicate(ctx context.Context, email string) (bool, error) {
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

func (r *UserRepository) GetByForgotPasswordToken(ctx context.Context, token, tokenResult string) (model.User, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user model.User

	filter := bson.M{"forgotPasswordToken.token": token, "forgotPasswordToken.resultToken": tokenResult, "forgotPasswordToken.tokenExpiresTime": bson.M{"$gt": time.Now()}}

	res := r.db.FindOne(nCtx, filter)

	if err := res.Decode(&user); err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, apperrors.ErrUserNotFound
		}
		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepository) ResetPassword(ctx context.Context, token, email, password string) error {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"forgotPasswordToken.token": token, "email": email}

	res, err := r.db.UpdateOne(nCtx, filter, bson.M{"$set": bson.M{"password": password}, "$unset": bson.M{"forgotPasswordToken": ""}})

	if res.MatchedCount == 0 {
		return apperrors.ErrUserNotFound
	}
	return err
}
