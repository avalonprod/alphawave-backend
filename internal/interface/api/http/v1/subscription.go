package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initSubscriptionRoutes(api *gin.RouterGroup) {
	subscription := api.Group("/subscription")
	{
		authenticated := subscription.Group("/", h.userIdentity)
		{
			teamSession := authenticated.Group("/", h.setTeamSessionFromCookie)
			{
				teamSession.POST("/create/", h.createSubscription)
			}

		}

	}

}

type createSubscriptionInput struct {
	PackageId string `json:"packageId"`
}

func (h *HandlerV1) createSubscription(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	teamId, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	var input createSubscriptionInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))

		return
	}

	if err := h.service.SubscriptionService.Create(c.Request.Context(), userId, input.PackageId, teamId); err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		if errors.Is(err, apperrors.ErrUserNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrUserNotFound.Error())

			return
		}
		logger.Errorf("failed to create subscription. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.Status(http.StatusCreated)
}
