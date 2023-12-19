package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initTasksRoutes(api *gin.RouterGroup) {
	tasks := api.Group("/tasks")
	{
		authenticated := tasks.Group("/", h.userIdentity, h.setTeamSessionFromCookie)
		{
			authenticated.POST("/create", h.createTask, h.checkRole(model.PERMISSION_ADD_NEW_TASKS))
			authenticated.GET("/:id", h.getByIdTask)
			authenticated.GET("", h.getAllTasks)
			authenticated.POST("/delete/:id", h.deleteTask)
			authenticated.POST("/finished/:id", h.finishedTask)
			authenticated.POST("/update", h.updateByIdTask, h.checkRole(model.PERMISSION_ADD_NEW_TASKS))
			authenticated.PUT("/update-position", h.updatePosition)
			authenticated.DELETE("/delete", h.deleteAll, h.checkRole(model.PERMISSION_ADD_NEW_TASKS))
			authenticated.DELETE("/clear", h.clearAll, h.checkRole(model.PERMISSION_ADD_NEW_TASKS))
			authenticated.PUT("/undo/:id", h.undoTask)
		}
	}
}

type CreateTaskInput struct {
	Title    string `json:"title"`
	Priority string `json:"priority"`
	Index    int    `json:"index"`
}

type CreateTaskResponse struct {
	Id string `json:"id"`
}

type UpdateTaskInput struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Index    int    `json:"index"`
}

func (h *HandlerV1) createTask(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}
	var input CreateTaskInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))

		return
	}

	id, err := h.service.TasksService.Create(c.Request.Context(), userID, types.TasksCreateDTO{
		Title:    input.Title,
		Priority: input.Priority,
		Index:    input.Index,
	})

	if err != nil {
		logger.Errorf("failed to create new task. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusCreated, CreateTaskResponse{Id: id})
}

func (h *HandlerV1) updatePosition(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}
	var input []types.UpdatePositionDTO
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))

		return
	}
	err = h.service.TasksService.UpdatePosition(c.Request.Context(), userID, input)
	if err != nil {
		logger.Errorf("failed to update tasks position. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.Status(http.StatusOK)
}

func (h *HandlerV1) getByIdTask(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "id is empty")

		return
	}

	task, err := h.service.TasksService.GetById(c.Request.Context(), userID, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *HandlerV1) getAllTasks(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}
	tasks, err := h.service.TasksService.GetAll(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *HandlerV1) updateByIdTask(c *gin.Context) {
	var input UpdateTaskInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))

		return
	}

	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	res, err := h.service.TasksService.UpdateById(c.Request.Context(), userID, types.UpdateTaskDTO{
		ID:       input.ID,
		Title:    input.Title,
		Status:   input.Status,
		Priority: input.Priority,
		Index:    input.Index,
	})

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusFound, res)
}

func (h *HandlerV1) deleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "id is empty")

		return
	}

	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	if err := h.service.TasksService.DeleteTaskById(c.Request.Context(), userID, id); err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.Status(http.StatusOK)
}

func (h *HandlerV1) finishedTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "id is empty")

		return
	}

	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	if err := h.service.TasksService.FinishedTaskById(c.Request.Context(), userID, id); err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}
	c.Status(http.StatusOK)
}

func (h *HandlerV1) deleteAll(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	err = h.service.TasksService.DeleteAll(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.Status(http.StatusOK)
}

func (h *HandlerV1) clearAll(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	err = h.service.TasksService.ClearAll(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.Status(http.StatusOK)
}

func (h *HandlerV1) undoTask(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "id is empty")

		return
	}
	if err := h.service.TasksService.UndoTask(c.Request.Context(), userID, id); err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		logger.Errorf("failed to undo task. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.Status(http.StatusOK)
}
