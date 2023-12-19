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

type PackagesRepository struct {
	db *mongo.Collection
}

func NewPackagesRepository(db *mongo.Database) *PackagesRepository {
	return &PackagesRepository{
		db: db.Collection(packagesCollection),
	}
}

func (r *PackagesRepository) CreateDefaultPackages(packages []model.Package) error {
	nCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := r.db.InsertOne(nCtx, packages); err != nil {
		return err
	}
	return nil
}

func (r *PackagesRepository) GetAllPackages(ctx context.Context) ([]model.Package, error) {
	nCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var packages []model.Package

	cur, err := r.db.Find(nCtx, bson.M{})
	if err != nil {
		return []model.Package{}, err
	}

	err = cur.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []model.Package{}, apperrors.ErrDocumentNotFound
		}
		return []model.Package{}, err
	}

	if err := cur.All(nCtx, &packages); err != nil {
		return []model.Package{}, err
	}
	return packages, nil
}

func (r *PackagesRepository) GetPackageById(ctx context.Context, packageId string) (model.Package, error) {
	nCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	fmt.Println(packageId)
	var pkg model.Package

	ObjectID, err := primitive.ObjectIDFromHex(packageId)
	if err != nil {
		return model.Package{}, err
	}
	filter := bson.M{"_id": ObjectID}

	res := r.db.FindOne(nCtx, filter)

	err = res.Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Package{}, apperrors.ErrDocumentNotFound
		}
		return model.Package{}, err
	}

	if err := res.Decode(&pkg); err != nil {
		return model.Package{}, err
	}

	return pkg, nil
}
