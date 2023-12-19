package v1

import (
	"errors"
	"net/http"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initPackagesRoutes(api *gin.RouterGroup) {
	packages := api.Group("/packages")
	{
		packages.GET("/get-all", h.getAll)

	}

}

func (h *HandlerV1) getAll(c *gin.Context) {

	packages, err := h.service.PackagesService.GetAll(c.Request.Context())

	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusOK, packages)
}
