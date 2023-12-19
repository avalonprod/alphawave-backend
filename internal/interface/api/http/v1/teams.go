package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
	"github.com/gin-gonic/gin"
)

func (h *HandlerV1) initTeamsRoutes(api *gin.RouterGroup) {
	teams := api.Group("/teams")
	{
		authenticated := teams.Group("/", h.userIdentity)
		{
			authenticated.POST("/create", h.createTeam)
			authenticated.GET("/my-own", h.getTeamByOwnerId)
			authenticated.GET("/set-session/:id", h.setSession)
			authenticated.GET("", h.getTeams)
			authenticated.GET("/:id")
			teamSession := authenticated.Group("/", h.setTeamSessionFromCookie)
			{
				teamSession.PUT("/settings", h.updateSettings)
				roles := teamSession.Group("/roles")
				{
					roles.GET("", h.getRoles)
					roles.PUT("", h.UpdatePermissions)
				}
			}
		}

	}

}

type createTeamInput struct {
	TeamName string `json:"teamName"`
	JobTitle string `json:"jobTitle"`
}

type createTeamResponse struct {
	Id string `json:"id"`
}

type updateTeamSettingsInput struct {
	UserActivityIndicator *bool   `json:"userActivityIndicator"`
	DisplayLinkPreview    *bool   `json:"displayLinkPreview"`
	DisplayFilePreview    *bool   `json:"displayFilePreview"`
	EnableGifs            *bool   `json:"enableGifs"`
	ShowWeekends          *bool   `json:"showWeekends"`
	FirstDayOfWeek        *string `json:"firstDayOfWeek"`
}

type teamSession struct {
	TeamID string
	Roles  []string
}

type getTeamResponse struct {
	ID       string `json:"id"`
	TeamName string `json:"teamName"`
	JobTitle string `json:"jobTitle"`
	OwnerID  string `json:"ownerID"`
}

type updatePermissionsInput struct {
	Role        string            `json:"role"`
	Permissions model.Permissions `jons:"permissions"`
}

func (h *HandlerV1) createTeam(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	var input createTeamInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))

		return
	}

	teamId, err := h.service.TeamsService.Create(c.Request.Context(), userID, types.CreateTeamsDTO{
		TeamName: input.TeamName,
		JobTitle: input.JobTitle,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			logger.Errorf("failed to create team. err: %v", err)
			newResponse(c, http.StatusNotFound, apperrors.ErrUserNotFound.Error())

			return
		}
		logger.Errorf("failed to create team. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	sessionData := teamSession{
		TeamID: teamId,
		Roles:  []string{},
	}

	member, err := h.service.MemberService.GetMemberByTeamIdAndUserId(c.Request.Context(), teamId, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrMemberNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrMemberNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	sessionData.Roles = member.Roles

	sessionDataJson, err := json.Marshal(sessionData)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "error: error marshal data to json")

		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "team_session",
		Value:    url.QueryEscape(string(sessionDataJson)),
		Path:     "/",
		Domain:   "",
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	c.JSON(http.StatusCreated, createTeamResponse{
		Id: teamId,
	})
}

func (h *HandlerV1) setSession(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		logger.Errorf("failed to set session for team. err: %v", err)
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	id := c.Param("id")
	if id == "" {
		newResponse(c, http.StatusBadRequest, "id is empty")

		return
	}

	team, err := h.service.TeamsService.GetTeamByID(c.Request.Context(), userID, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrTeamNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrTeamNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	sessionData := teamSession{
		TeamID: team.ID,
		Roles:  []string{},
	}

	member, err := h.service.MemberService.GetMemberByTeamIdAndUserId(c.Request.Context(), team.ID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrMemberNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrMemberNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}
	sessionData.Roles = member.Roles

	sessionDataJson, err := json.Marshal(sessionData)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, "error: error marshal data to json")

		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "team_session",
		Value:    url.QueryEscape(string(sessionDataJson)),
		Path:     "/",
		Domain:   "",
		MaxAge:   0,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
}

func (h *HandlerV1) getTeams(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}
	teams, err := h.service.TeamsService.GetTeamsByUser(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrTeamNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrTeamNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusOK, teams)
}

// TODO
func (h *HandlerV1) getTeamById(c *gin.Context) {

}

func (h *HandlerV1) updateSettings(c *gin.Context) {
	id, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	var input updateTeamSettingsInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))

		return
	}

	err = h.service.TeamsService.UpdateTeamSettings(c.Request.Context(), id, types.UpdateTeamSettingsDTO{
		UserActivityIndicator: input.UserActivityIndicator,
		DisplayLinkPreview:    input.DisplayLinkPreview,
		DisplayFilePreview:    input.DisplayFilePreview,
		EnableGifs:            input.EnableGifs,
		ShowWeekends:          input.ShowWeekends,
		FirstDayOfWeek:        input.FirstDayOfWeek,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrTeamNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrTeamNotFound.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.Status(http.StatusOK)
}

func (h *HandlerV1) getTeamByOwnerId(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	team, err := h.service.TeamsService.GetTeamByOwnerId(c.Request.Context(), userId)
	if err != nil {
		if errors.Is(err, apperrors.ErrTeamNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrTeamNotFound.Error())

			return
		}
		if errors.Is(err, apperrors.ErrInvalidIdFormat) {
			newResponse(c, http.StatusBadRequest, apperrors.ErrInvalidIdFormat.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusOK, getTeamResponse{
		ID:       team.ID,
		TeamName: team.TeamName,
		JobTitle: team.JobTitle,
		OwnerID:  team.OwnerID,
	})
}

func (h *HandlerV1) getRoles(c *gin.Context) {
	id, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	roles, err := h.service.RolesService.GetRolesByTeamId(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrDocumentNotFound) {
			newResponse(c, http.StatusNotFound, apperrors.ErrDocumentNotFound.Error())

			return
		}
		if errors.Is(err, apperrors.ErrInvalidIdFormat) {
			newResponse(c, http.StatusBadRequest, apperrors.ErrInvalidIdFormat.Error())

			return
		}

		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	c.JSON(http.StatusOK, roles)
}

func (h *HandlerV1) UpdatePermissions(c *gin.Context) {
	id, err := getTeamId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, apperrors.ErrInternalServerError.Error())

		return
	}

	var input []updatePermissionsInput

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, fmt.Sprintf("Incorrect data format. err: %v", err))

		return
	}

	roles := make([]types.UpdatePermissionsDTO, len(input))

	for i := range roles {
		roles[i] = types.UpdatePermissionsDTO{
			Role:        input[i].Role,
			Permissions: input[i].Permissions,
		}
	}

	err = h.service.RolesService.UpdatePermissions(c.Request.Context(), id, roles)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidIdFormat) {
			newResponse(c, http.StatusBadRequest, apperrors.ErrInvalidIdFormat.Error())

			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.Status(http.StatusOK)
}
