package handlers

import (
	"net/http"

	"clairvoyance/internal/services"
	"clairvoyance/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler struct {
	teamService *services.TeamService
	credService *services.CredentialService
}

func NewTeamHandler(teamService *services.TeamService, credService *services.CredentialService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
		credService: credService,
	}
}

// GetTeam handles GET /api/teams/:id
// Public. Accepts a UUID or a slug.
func (th *TeamHandler) GetTeam(c *gin.Context) {
	idParam := c.Param("id")

	// Try as UUID first, fall back to slug
	teamID, err := uuid.Parse(idParam)
	if err != nil {
		team, err := th.teamService.GetTeamBySlug(idParam)
		if err != nil {
			if err == utils.ErrNotFound {
				utils.ErrorResponse(c, http.StatusNotFound, "team not found")
				return
			}
			utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SuccessResponse(c, http.StatusOK, team)
		return
	}

	team, err := th.teamService.GetTeamByID(teamID)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "team not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, team)
}

// GetTeamCredentials handles GET /api/teams/:id/credentials
// Public. Returns all credentials held by any member of the team.
func (th *TeamHandler) GetTeamCredentials(c *gin.Context) {
	idParam := c.Param("id")
	teamID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid team id")
		return
	}

	// Collect all member email/wallet refs
	refs, err := th.teamService.GetMemberRefs(teamID)
	if err != nil {
		if err == utils.ErrNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, "team not found")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(refs) == 0 {
		utils.SuccessResponse(c, http.StatusOK, gin.H{"credentials": []any{}})
		return
	}

	// Fetch all credentials matching those refs
	credentials, err := th.credService.GetCredentialsByRefs(refs)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"credentials": credentials})
}
