package v1

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	maxUploadSize = 2 << 30
)

// type contentRange struct {
// 	rangeStart int64
// 	rangeEnd   int64
// 	fileSize   int64
// }

type CreateFolderInput struct {
	Name         string `json:"name"`
	ParentFolder string `json:"parentFolder"`
}

type RenameFileInput struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type RenameFileResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type FileMetadataResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Size      int       `json:"size"`
	Extension string    `json:"extension"`
	File      http.File `json:"file"`
}

func (h *HandlerV1) initFilesRoutes(api *gin.RouterGroup) {
	files := api.Group("/files")
	{
		authenticated := files.Group("/", h.userIdentity)
		{
			teamSession := authenticated.Group("/", h.setTeamSessionFromCookie)
			{
				teamSession.GET("/url/:id", h.getFilePresignedURL)
				teamSession.DELETE("/:id", h.deleteFile)
				teamSession.POST("", h.createFile)
				teamSession.PATCH("", h.renameFile)
				folders := teamSession.Group("/folders")
				{
					folders.GET("/root", h.getRootFolder)
					folders.POST("", h.createFolder)
					folders.GET("/:id", h.getFolderContent)
				}
			}
		}

	}
}

func (h *HandlerV1) createFile(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	teamID, err := getTeamId(c)
	if err != nil {
		logger.Errorf("failed to get team id. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	userID, err := getUserId(c)
	if err != nil {
		logger.Errorf("failed to get team user id. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	tempFile, err := c.FormFile("file")

	if err != nil {
		logger.Error(err)
		newResponse(c, http.StatusBadRequest, err.Error())

		return
	}

	folder := c.PostForm("folder")
	fileName := c.PostForm("fileName")

	fileReader, err := tempFile.Open()

	if err != nil {
		logger.Error(err)
		newResponse(c, http.StatusBadRequest, fmt.Errorf("failed to open file").Error())

		return
	}

	defer fileReader.Close()

	fileData, err := ioutil.ReadAll(fileReader)

	if err != nil {
		logger.Error(err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	reader := bytes.NewReader(fileData)

	file, err := h.service.FilesService.Create(c.Request.Context(), reader, types.CreateFileDTO{
		TeamID:    teamID,
		UserID:    userID,
		FileName:  fileName,
		Type:      filepath.Ext(tempFile.Filename),
		Extension: filepath.Ext(tempFile.Filename),
		Size:      len(fileData),
		Folder:    folder,
	})

	if err != nil {
		logger.Errorf("failed to upload file. err: %v", err)
		newResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to upload file").Error())

		return
	}

	c.JSON(http.StatusCreated, file)
}

func (h *HandlerV1) renameFile(c *gin.Context) {
	teamID, err := getTeamId(c)
	if err != nil {
		logger.Errorf("failed to get team id. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	var input RenameFileInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))

		return
	}

	fileName, err := h.service.FilesService.RenameFile(c.Request.Context(), teamID, input.Id, input.Name)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusOK, RenameFileResponse{
		Id:   input.Id,
		Name: fileName,
	})
}

func (h *HandlerV1) createFolder(c *gin.Context) {
	teamID, err := getTeamId(c)
	if err != nil {
		logger.Errorf("failed to get team id. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	var input CreateFolderInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))

		return
	}

	folder, err := h.service.FilesService.CreateFolder(c.Request.Context(), types.CreateFolderDTO{
		TeamID:         teamID,
		FolderName:     input.Name,
		ParentFolderId: input.ParentFolder,
	})
	if err != nil {
		logger.Errorf("failed to create folder. err: %v", err)
		newResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to create folder. folder name: %s", input.Name).Error())

		return
	}

	c.JSON(http.StatusCreated, folder)
}

func (h *HandlerV1) getFilePresignedURL(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		logger.Warnf("failed to get file id. file id is empty")
		newResponse(c, http.StatusBadRequest, "id is empty")

		return
	}

	teamID, err := getTeamId(c)
	if err != nil {
		logger.Errorf("failed to get team id. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	url, err := h.service.FilesService.GetFilePresignedURL(c.Request.Context(), teamID, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			logger.Warnf("failed to get file presigned url. err: %v", err)
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		logger.Errorf("failed to get file presigned url. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("https://%s", url))
}

func (h *HandlerV1) getFolderContent(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "id is empty")

		return
	}

	teamID, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	content, err := h.service.FilesService.GetFolderContent(c.Request.Context(), teamID, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			logger.Warnf("failed to get folder content. err: %v", err)
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		logger.Errorf("failed to get folder content. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusOK, content)
}

func (h *HandlerV1) getRootFolder(c *gin.Context) {
	teamID, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}
	folder, err := h.service.FilesService.GetFolderRoot(c.Request.Context(), teamID)

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			logger.Warnf("failed to get folder content. err: %v", err)
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		logger.Errorf("failed to get folder content. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusOK, folder)
}

func (h *HandlerV1) deleteFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "id is empty")

		return
	}

	teamID, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	if err := h.service.FilesService.Delete(c.Request.Context(), teamID, id); err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			logger.Warnf("failed to delete file. err: %v", err)
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		logger.Errorf("failed to delete file. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.Status(http.StatusNoContent)
}
