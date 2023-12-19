package repository

import (
	"context"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type PackagesRepository interface {
	GetPackageById(ctx context.Context, packageId string) (model.Package, error)
	CreateDefaultPackages(packages []model.Package) error
	GetAllPackages(ctx context.Context) ([]model.Package, error)
}
