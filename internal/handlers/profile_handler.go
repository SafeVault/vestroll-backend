package handlers

import (
    "net/http"

    "github.com/codeZe-us/vestroll-backend/internal/models"
    "github.com/codeZe-us/vestroll-backend/internal/services"
    "github.com/gin-gonic/gin"
)

// ProfileHandler manages onboarding profile endpoints (account type, personal, address)
type ProfileHandler struct {
    service *services.ProfileService
}

func NewProfileHandler(s *services.ProfileService) *ProfileHandler { return &ProfileHandler{service: s} }

// RegisterRoutes adds routes under /api(/v1)/profile
func (h *ProfileHandler) RegisterRoutes(router *gin.RouterGroup) {
    router.POST("/account-type", h.PostAccountType)
    router.POST("/personal-details", h.PostPersonalDetails)
    router.POST("/address", h.PostAddress)
    router.GET("/status", h.GetStatus)
}

// POST /api/profile/account-type
func (h *ProfileHandler) PostAccountType(c *gin.Context) {
    var req models.AccountTypeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
        return
    }
    prof, err := h.service.UpdateAccountType(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
        return
    }
    c.JSON(http.StatusOK, models.ProfileResponse{Success: true, Message: "Account type saved", Profile: prof})
}

// POST /api/profile/personal-details
func (h *ProfileHandler) PostPersonalDetails(c *gin.Context) {
    var req models.PersonalDetailsRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
        return
    }
    prof, err := h.service.UpdatePersonalDetails(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
        return
    }
    c.JSON(http.StatusOK, models.ProfileResponse{Success: true, Message: "Personal details saved", Profile: prof})
}

// POST /api/profile/address
func (h *ProfileHandler) PostAddress(c *gin.Context) {
    var req models.AddressRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
        return
    }
    prof, err := h.service.UpdateAddress(c.Request.Context(), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: err.Error()})
        return
    }
    c.JSON(http.StatusOK, models.ProfileResponse{Success: true, Message: "Address saved", Profile: prof})
}

// GET /api/profile/status?user_id=...
func (h *ProfileHandler) GetStatus(c *gin.Context) {
    userID := c.Query("user_id")
    if userID == "" {
        c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "validation_error", Message: "user_id is required"})
        return
    }
    prof, err := h.service.GetProfile(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "internal_error", Message: err.Error()})
        return
    }
    c.JSON(http.StatusOK, models.ProfileResponse{Success: true, Message: "Profile status", Profile: prof})
}
