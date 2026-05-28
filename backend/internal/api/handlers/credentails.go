package handlers

import (
	"net/http"

	"clairvoyance/internal/services"
	"clairvoyance/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CredentialHandler struct {
	credService *services.CredentialService
}

func NewCredentialHandler(credService *services.CredentialService) *CredentialHandler {
	return &CredentialHandler{credService: credService}
}

// Issue handles POST /api/credentials/issue
// Requires issuer or leader role
func (ch *CredentialHandler) Issue(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	issuerUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "invalid user id")
		return
	}

	var req services.IssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := ch.credService.IssueCredential(issuerUUID, req)
	if err != nil {
		if err == utils.ErrForbidden {
			utils.ErrorResponse(c, http.StatusForbidden, "insufficient permissions")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, resp)
}

// Verify handles GET /api/credentials/verify/:id
// Public endpoint - no authentication required
func (ch *CredentialHandler) Verify(c *gin.Context) {
	idParam := c.Param("id")
	credentialID, err := uuid.Parse(idParam)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid credential id")
		return
	}

	resp, err := ch.credService.VerifyCredential(credentialID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Return appropriate HTTP status based on verification result
	switch resp.Status {
	case "VERIFIED":
		utils.SuccessResponse(c, http.StatusOK, resp)
	case "NOT_FOUND":
		c.JSON(http.StatusNotFound, resp)
	case "TAMPERED", "ERROR":
		utils.SuccessResponse(c, http.StatusOK, resp) // Still 200 but with status field
	default:
		utils.SuccessResponse(c, http.StatusOK, resp)
	}
}

// GetMine handles GET /api/credentials/mine
// Returns credentials owned by the authenticated user
func (ch *CredentialHandler) GetMine(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "invalid user id")
		return
	}

	credentials, err := ch.credService.GetUserCredentials(userUUID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"credentials": credentials,
	})
}

// GetIssued handles GET /api/credentials/issued
// Returns credentials issued by the authenticated user (requires issuer role)
func (ch *CredentialHandler) GetIssued(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		utils.ErrorResponse(c, http.StatusInternalServerError, "invalid user id")
		return
	}

	credentials, err := ch.credService.GetIssuerCredentials(userUUID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"credentials": credentials,
	})
}