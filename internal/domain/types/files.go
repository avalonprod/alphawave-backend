package types

import (
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
)

type CreateFileDTO struct {
	UserID    string `json:"userId"`
	TeamID    string `json:"teamId"`
	FileName  string `json:"fileName"`
	Type      string `json:"type"`
	Extension string `json:"extension"`
	Size      int    `json:"size"`
	Folder    string `json:"folder"`
}

type UploadImageDTO struct {
	FileName  string `json:"fileName"`
	Extension string `json:"extension"`
	Size      int    `json:"size"`
}

type GetFileDTO struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Size      int     `json:"size"`
	Extension string  `json:"extention"`
	File      *[]byte `json:"file"`
}

type CreateFolderDTO struct {
	TeamID         string `json:"teamId"`
	FolderName     string `json:"folderName"`
	ParentFolderId string `json:"parentFolderId"`
}

type FolderDTO struct {
	Id               string             `json:"id"`
	Name             string             `json:"name"`
	Type             string             `json:"type"`
	Path             []model.FolderPath `json:"path"`
	CreatedAt        time.Time          `json:"createdAt"`
	LastModifiedTime time.Time          `json:"lastModifiedTime"`
	ParentFolderId   string             `json:"parentFolderId"`
}

type FileDTO struct {
	ID        string    `json:"id"`
	OwnerName string    `json:"ownerName"`
	Name      string    `json:"name"`
	FolderId  string    `json:"folderId"`
	Url       string    `json:"url"`
	Type      string    `json:"type"`
	Path      []string  `json:"path"`
	Size      int       `json:"size"`
	Extension string    `json:"extension"`
	CreatedAt time.Time `json:"createdAt"`
}

type FolderContentDTO struct {
	FolderInfo FolderDTO      `json:"folderInfo"`
	Folders    []model.Folder `json:"folders"`
	Files      []model.File   `json:"files"`
}

type RootFolderInfoDTO struct {
	ID               string             `json:"id"`
	Name             string             `json:"name"`
	Type             string             `json:"type"`
	Path             []model.FolderPath `json:"path"`
	CreatedAt        time.Time          `json:"createdAt"`
	LastModifiedTime time.Time          `json:"lastModifiedTime"`
}

type RootFolderDTO struct {
	FolderInfo RootFolderInfoDTO `json:"folderInfo"`
	Folders    []FolderDTO       `json:"folders"`
	Files      []FileDTO         `json:"files"`
}

// IMAGES
type ImageOutputDTO struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}
