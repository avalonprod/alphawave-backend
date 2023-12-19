package repository

import (
	"context"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type FilesRepository interface {
	Create(ctx context.Context, input model.File) (string, error)
	RenameFile(ctx context.Context, teamId, fileId, name string, path []string) error
	Delete(ctx context.Context, teamID, fileID string) error
	GetFileById(ctx context.Context, teamID, fileID string) (model.File, error)
	GetFilesByFolderId(ctx context.Context, teamId, folderId string) ([]model.File, error)
}
