package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/codegenerator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type storageProvider interface {
	UploadFile(ctx context.Context, bucketName, objectName, fileName string, fileSize int64, reader io.Reader) error
	// CreateFolder(ctx context.Context, bucketName, objectName string) error
	GetFilePresignedURL(ctx context.Context, bucketName, fileName string, expiresTime time.Duration) (string, error)
	GetFile(ctx context.Context, bucketName, fileName string) (*[]byte, error)
	DeleteFile(ctx context.Context, bucketName, fileName string) error
}

const (
	IMAGE = "image"
	VIDEO = "video"
	OTHER = "other"
)

type FilesService struct {
	storageProvider  storageProvider
	repository       repository.FilesRepository
	folderRepository repository.FolderRepository
	codeGenerator    *codegenerator.CodeGenerator
	userRepository   repository.UserRepository

	endpointURL string
}

const defaultFileExpiresTime = time.Second * 60 * 60 * 24
const FOLDER_TYPE_ROOT = "root"
const FOLDER_TYPE_DEFAULT = "default"

func NewFilesService(storageProvider storageProvider, repository repository.FilesRepository, codeGenerator *codegenerator.CodeGenerator, folderRepository repository.FolderRepository, endpointURL string, userRepository repository.UserRepository) *FilesService {
	return &FilesService{
		storageProvider:  storageProvider,
		repository:       repository,
		folderRepository: folderRepository,
		codeGenerator:    codeGenerator,
		userRepository:   userRepository,
		endpointURL:      endpointURL,
	}
}

func (s *FilesService) Create(ctx context.Context, reader io.Reader, input types.CreateFileDTO) (types.FileDTO, error) {
	uuid := s.codeGenerator.GenerateUUID()

	fileName := fmt.Sprintf("%s%s", uuid, input.Extension)
	var folderId string
	var path []string
	if strings.Trim(input.Folder, " ") == "" {
		rootFolder, err := s.GetFolderRoot(ctx, input.TeamID)
		if err != nil {
			return types.FileDTO{}, err
		}

		folderId = rootFolder.FolderInfo.ID
		path = append(path, rootFolder.FolderInfo.Name, input.FileName)
	} else {
		folder, err := s.folderRepository.GetFolderById(ctx, input.TeamID, input.Folder)
		if err != nil {
			return types.FileDTO{}, err
		}

		folderId = folder.ID
		for _, item := range folder.Path {
			path = append(path, item.Name)
		}

		path = append(path, folder.Name, input.FileName)
	}
	if strings.Trim(input.FileName, " ") == "" {
		return types.FileDTO{}, apperrors.ErrFileNameIsEmpty
	}

	if err := s.storageProvider.UploadFile(ctx, model.BUCKET_STORAGE, fileName, input.FileName, int64(input.Size), reader); err != nil {
		return types.FileDTO{}, err
	}

	user, err := s.userRepository.GetUserById(ctx, input.UserID)
	if err != nil {
		return types.FileDTO{}, err
	}

	var fileType string

	if input.Type == ".jpg" || input.Type == ".jpeg" || input.Type == ".png" {
		fileType = IMAGE
	} else {
		fileType = OTHER
	}

	createFile := model.File{
		TeamId:    input.TeamID,
		Name:      input.FileName,
		OwnerName: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		FilePath:  fileName,
		Url:       s.generateFileURL(fileName),
		Key:       uuid,
		Type:      fmt.Sprintf("%s/%s", fileType, input.Type[1:]),
		Path:      path,
		Size:      input.Size,
		Extension: input.Extension[1:],
		CreatedAt: time.Now(),
		FolderId:  folderId,
	}

	id, err := s.repository.Create(ctx, createFile)
	if err != nil {
		return types.FileDTO{}, err
	}
	return types.FileDTO{
		ID:        id,
		OwnerName: createFile.OwnerName,
		Name:      createFile.Name,
		FolderId:  createFile.FolderId,
		Url:       createFile.Url,
		Type:      createFile.Type,
		Path:      createFile.Path,
		Size:      createFile.Size,
		Extension: createFile.Extension,
		CreatedAt: createFile.CreatedAt,
	}, nil
}

func (s *FilesService) RenameFile(ctx context.Context, teamId string, fileId string, name string) (string, error) {
	file, err := s.repository.GetFileById(ctx, teamId, fileId)
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s.%s", name, file.Extension)
	path := file.Path
	path[len(path)-1] = fileName
	err = s.repository.RenameFile(ctx, teamId, fileId, fileName, path)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

func (s *FilesService) CreateRootFolder(ctx context.Context, teamId string) error {
	folderId := primitive.NewObjectID().Hex()
	path := make([]model.FolderPath, 1)
	path[0] = model.FolderPath{
		Id:   folderId,
		Name: "root",
	}

	if _, err := s.folderRepository.CreateFolder(ctx, teamId, model.Folder{
		ID:               folderId,
		TeamId:           teamId,
		Name:             FOLDER_TYPE_ROOT,
		Type:             FOLDER_TYPE_ROOT,
		Path:             path,
		CreatedAt:        time.Now(),
		LastModifiedTime: time.Now(),
	}); err != nil {
		return err
	}
	return nil
}

func (s *FilesService) GetFolderRoot(ctx context.Context, teamId string) (types.RootFolderDTO, error) {
	folder, err := s.folderRepository.GetFolderRoot(ctx, teamId, FOLDER_TYPE_ROOT)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return types.RootFolderDTO{}, apperrors.ErrDocumentNotFound
		}
		return types.RootFolderDTO{}, err
	}

	foldersContent, err := s.folderRepository.GetFolderContentById(ctx, teamId, folder.ID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return types.RootFolderDTO{}, apperrors.ErrDocumentNotFound
		}
		return types.RootFolderDTO{}, err
	}
	filesContent, err := s.repository.GetFilesByFolderId(ctx, teamId, folder.ID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return types.RootFolderDTO{}, apperrors.ErrDocumentNotFound
		}
		return types.RootFolderDTO{}, err
	}
	var wg sync.WaitGroup
	folders := make([]types.FolderDTO, len(foldersContent))
	files := make([]types.FileDTO, len(filesContent))
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i, item := range foldersContent {
			folders[i] = types.FolderDTO{
				Id:               item.ID,
				Name:             item.Name,
				Type:             item.Type,
				CreatedAt:        item.CreatedAt,
				Path:             item.Path,
				LastModifiedTime: item.LastModifiedTime,
				ParentFolderId:   item.ParentFolder,
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i, item := range filesContent {
			files[i] = types.FileDTO{
				ID:        item.ID,
				OwnerName: item.OwnerName,
				Name:      item.Name,
				FolderId:  item.FolderId,
				Url:       item.Url,
				Type:      item.Type,
				Path:      item.Path,
				Size:      item.Size,
				Extension: item.Extension,
				CreatedAt: item.CreatedAt,
			}
		}
	}()

	wg.Wait()
	return types.RootFolderDTO{
		FolderInfo: types.RootFolderInfoDTO{
			ID:               folder.ID,
			Name:             folder.Name,
			Type:             folder.Type,
			Path:             folder.Path,
			CreatedAt:        folder.CreatedAt,
			LastModifiedTime: folder.LastModifiedTime,
		},
		Folders: folders,
		Files:   files,
	}, nil
}

func (s *FilesService) CreateFolder(ctx context.Context, input types.CreateFolderDTO) (types.FolderDTO, error) {
	var parentFolderId string
	var path []model.FolderPath
	folderId := primitive.NewObjectID().Hex()

	if strings.Trim(input.ParentFolderId, " ") == "" {
		rootFolder, err := s.GetFolderRoot(ctx, input.TeamID)
		if err != nil {
			return types.FolderDTO{}, err
		}
		parentFolderId = rootFolder.FolderInfo.ID
		path = append(path, model.FolderPath{
			Id:   rootFolder.FolderInfo.ID,
			Name: rootFolder.FolderInfo.Name,
		})
	} else {
		folder, err := s.folderRepository.GetFolderById(ctx, input.TeamID, input.ParentFolderId)
		if err != nil {
			return types.FolderDTO{}, err
		}

		parentFolderId = folder.ID
		path = folder.Path
		path = append(path, model.FolderPath{
			Id:   folder.ID,
			Name: folder.Name,
		})
	}
	createFolder := model.Folder{
		ID:               folderId,
		TeamId:           input.TeamID,
		Name:             input.FolderName,
		Type:             FOLDER_TYPE_DEFAULT,
		Path:             path,
		ParentFolder:     parentFolderId,
		LastModifiedTime: time.Now(),
		CreatedAt:        time.Now(),
	}
	id, err := s.folderRepository.CreateFolder(ctx, input.TeamID, createFolder)
	if err != nil {
		return types.FolderDTO{}, err
	}
	return types.FolderDTO{
		Id:               id,
		Name:             createFolder.Name,
		Type:             createFolder.Type,
		Path:             createFolder.Path,
		CreatedAt:        createFolder.CreatedAt,
		LastModifiedTime: createFolder.LastModifiedTime,
		ParentFolderId:   createFolder.ParentFolder,
	}, nil
}

func (s *FilesService) GetFolderContent(ctx context.Context, teamId, folderID string) (types.FolderContentDTO, error) {
	folder, err := s.folderRepository.GetFolderById(ctx, teamId, folderID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return types.FolderContentDTO{}, apperrors.ErrDocumentNotFound
		}
		return types.FolderContentDTO{}, err
	}
	folders, err := s.folderRepository.GetFolderContentById(ctx, teamId, folderID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return types.FolderContentDTO{}, apperrors.ErrDocumentNotFound
		}
		return types.FolderContentDTO{}, err
	}
	files, err := s.repository.GetFilesByFolderId(ctx, teamId, folderID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return types.FolderContentDTO{}, apperrors.ErrDocumentNotFound
		}
		return types.FolderContentDTO{}, err
	}

	return types.FolderContentDTO{
		FolderInfo: types.FolderDTO{
			Id:               folder.ID,
			Name:             folder.Name,
			Type:             folder.Type,
			Path:             folder.Path,
			CreatedAt:        folder.CreatedAt,
			LastModifiedTime: folder.LastModifiedTime,
			ParentFolderId:   folder.ParentFolder,
		},
		Folders: folders,
		Files:   files,
	}, nil
}

func (s *FilesService) GetFilePresignedURL(ctx context.Context, teamID, fileID string) (string, error) {
	file, err := s.repository.GetFileById(ctx, teamID, fileID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return "", apperrors.ErrDocumentNotFound
		}
		return "", err
	}

	url, err := s.storageProvider.GetFilePresignedURL(ctx, model.BUCKET_STORAGE, file.FilePath, defaultFileExpiresTime)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *FilesService) Delete(ctx context.Context, teamId, fileId string) error {
	res, err := s.repository.GetFileById(ctx, teamId, fileId)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			return apperrors.ErrDocumentNotFound
		}
		return err
	}

	if err := s.storageProvider.DeleteFile(ctx, model.BUCKET_STORAGE, res.FilePath); err != nil {
		return err
	}
	if err := s.repository.Delete(ctx, teamId, res.ID); err != nil {
		return err
	}
	return nil
}

func (s *FilesService) UploadImage(ctx context.Context, reader io.Reader, input types.UploadImageDTO) (types.ImageOutputDTO, error) {
	uuid := s.codeGenerator.GenerateUUID()
	fileName := fmt.Sprintf("%s%s", uuid, input.Extension)

	if input.Extension != ".jpg" && input.Extension != ".svg" && input.Extension != ".png" && input.Extension != ".jpeg" {
		return types.ImageOutputDTO{}, apperrors.ErrInvalidFileType
	}

	if err := s.storageProvider.UploadFile(ctx, model.BUCKET_IMAGES, fileName, input.FileName, int64(input.Size), reader); err != nil {
		return types.ImageOutputDTO{}, err
	}
	return types.ImageOutputDTO{
		Url:  s.generateImageURL(fileName),
		Path: fileName,
	}, nil
}

func (s *FilesService) DeleteImage(ctx context.Context, filePath string) error {
	if err := s.storageProvider.DeleteFile(ctx, model.BUCKET_IMAGES, filePath); err != nil {
		return err
	}
	return nil
}

func (s *FilesService) generateFileURL(filename string) string {
	return fmt.Sprintf("https://%s/%s/%s", s.endpointURL, model.BUCKET_STORAGE, filename)
}
func (s *FilesService) generateImageURL(filename string) string {
	return fmt.Sprintf("https://%s/%s/%s", s.endpointURL, model.BUCKET_IMAGES, filename)
}
