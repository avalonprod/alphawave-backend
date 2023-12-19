package service

import (
	"context"
	"errors"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
)

type PackagesService struct {
	repository repository.PackagesRepository
}

func NewPackagesService(repository repository.PackagesRepository) *PackagesService {
	return &PackagesService{
		repository: repository,
	}
}

// func (s *PackagesService) CreateDefaultPackages() error {
// 	var feature model.Feature

// 	defaultPackage := model.Package{
// 		Name:        "Start",
// 		Description: "Package for start",
// 		Price:       20,
// 		Currency:    "usd",
// 	}
// }

func (s *PackagesService) GetAll(ctx context.Context) ([]model.Package, error) {
	packages, err := s.repository.GetAllPackages(ctx)

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return []model.Package{}, apperrors.ErrDocumentNotFound
		}
		return []model.Package{}, err
	}

	return packages, nil
}

func (s *PackagesService) GetById(ctx context.Context, packageId string) (model.Package, error) {
	pkg, err := s.repository.GetPackageById(ctx, packageId)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return model.Package{}, apperrors.ErrDocumentNotFound
		}
		return model.Package{}, err
	}
	return pkg, nil
}
