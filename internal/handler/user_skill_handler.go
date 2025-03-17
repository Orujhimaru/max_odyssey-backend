package handler

import (
	"encoding/json"
	"log"
	"max-odyssey-backend/internal/middleware"
	"max-odyssey-backend/internal/service"
	"max-odyssey-backend/utils"
	"net/http"
)

type UserSkillHandler struct {
	service *service.UserSkillService
}

func NewUserSkillHandler(service *service.UserSkillService) *UserSkillHandler {
	return &UserSkillHandler{service: service}
}

func (h *UserSkillHandler) GetUserSkills(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		log.Println("User not found in context")
		utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Get user skills
	skills, err := h.service.GetUserSkills(int64(user.ID))
	if err != nil {
		log.Printf("Error getting user skills: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch user skills", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"skills": skills,
	})
}

func (h *UserSkillHandler) CreateOrUpdateUserSkill(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		log.Println("User not found in context")
		utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Parse request body
	var requestBody struct {
		SkillName  string  `json:"skill_name"`
		SkillScore float64 `json:"skill_score"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate input
	if requestBody.SkillName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Skill name is required", nil)
		return
	}

	// Create or update skill
	skill, err := h.service.CreateOrUpdateUserSkill(int64(user.ID), requestBody.SkillName, requestBody.SkillScore)
	if err != nil {
		log.Printf("Error creating/updating skill: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create/update skill", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, skill)
}

func (h *UserSkillHandler) DeleteUserSkill(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		log.Println("User not found in context")
		utils.RespondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	// Parse request body
	var requestBody struct {
		SkillName string `json:"skill_name"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("Error parsing request body: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate input
	if requestBody.SkillName == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Skill name is required", nil)
		return
	}

	// Delete skill
	err = h.service.DeleteUserSkill(int64(user.ID), requestBody.SkillName)
	if err != nil {
		log.Printf("Error deleting skill: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete skill", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
