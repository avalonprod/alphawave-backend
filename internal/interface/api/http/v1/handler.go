package v1

import (
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/service"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/auth/manager"
	"github.com/gin-gonic/gin"
)

type HandlerV1 struct {
	service         *service.Service
	JWTManager      *manager.JWTManager
	refreshTokenTTL time.Duration
	frontEndUrl     string
}

func NewHandler(service *service.Service, JWTManager *manager.JWTManager, refreshTokenTTL time.Duration, frontEndUrl string) *HandlerV1 {
	return &HandlerV1{
		service:         service,
		JWTManager:      JWTManager,
		refreshTokenTTL: refreshTokenTTL,
		frontEndUrl:     frontEndUrl,
	}
}

func (h *HandlerV1) InitRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUserRoutes(v1)
		h.initTasksRoutes(v1)
		h.initTeamsRoutes(v1)
		h.initMembersRoutes(v1)
		h.initAiChatRoutes(v1)
		h.initFilesRoutes(v1)
		h.initPaymentRoutes(v1)
		h.initSubscriptionRoutes(v1)
		h.initPackagesRoutes(v1)
	}
}
