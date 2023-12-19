package v1

import (
	"errors"
	"net/http"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initPaymentRoutes(api *gin.RouterGroup) {
	payment := api.Group("/payment")
	{
		authenticated := payment.Group("/", h.userIdentity)
		{
			teamSession := authenticated.Group("/", h.setTeamSessionFromCookie)
			{
				teamSession.POST("/new-method", h.createPaymentMethod)
			}

		}

	}

}

type createPaymentMethodResponse struct {
	PaymentIntentId string `json:"paymentMethodId"`
}

func (h *HandlerV1) createPaymentMethod(c *gin.Context) {
	teamID, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	paymentIntentID, err := h.service.PaymentService.CreateNewPaymentMethod(c.Request.Context(), teamID)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		logger.Errorf("failed to create payment method. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusCreated, createPaymentMethodResponse{PaymentIntentId: *paymentIntentID})
}
