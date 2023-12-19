package repository

import (
	"context"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type FolderRepository interface {
	CreateFolder(ctx context.Context, teamId string, folder model.Folder) (string, error)
	GetFolderContentById(ctx context.Context, teamId, id string) ([]model.Folder, error)
	GetFolderRoot(ctx context.Context, teamId string, folderType string) (model.Folder, error)
	GetFolderById(ctx context.Context, teamId string, folderId string) (model.Folder, error)
}
