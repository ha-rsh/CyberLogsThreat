package handlers

import (
	"cybersecuritySystem/shared/logger"
	"cybersecuritySystem/shared/utils"
	"cybersecuritySystem/threat-analyzer-service/service"
	"net/http"

	"github.com/gorilla/mux"
)

type ThreatHandler struct {
	service 		*service.ThreatService
}

func NewThreatHandler(service *service.ThreatService) *ThreatHandler {
	return &ThreatHandler{service: service}
}

// AnalyzeThreats godoc
// @Summary Analyze threats
// @Description Run threat detection analysis on all logs in the system
// @Tags threats
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security ApiKeyAuth
// @Router /api/threats/analyze [post]
func (h *ThreatHandler) AnalyzeThreats(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.AnalyzeThreats()
	if err != nil {
		logger.Error("Threat analysis failed: %v", err)
		utils.SendErrorResponse(w, "Failed to analyze threats", http.StatusInternalServerError, err.Error())
		return
	}
	
	logger.Info("Threat analysis completed: %d threats detected", result.ThreatsDetected)
	utils.SendSuccessResponse(w, result, nil)
}

// GetThreats godoc
// @Summary Get all threats
// @Description Retrieve all detected threats from the system
// @Tags threats
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security ApiKeyAuth
// @Router /api/threats [get]
func (h *ThreatHandler) GetThreats(w http.ResponseWriter, r *http.Request) {
	threats, err := h.service.GetAllThreats()
	if err != nil {
		logger.Error("Failed to retrieve threats: %v", err)
		utils.SendErrorResponse(w, "Failed to retrieve threats", http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, threats, &utils.Meta{Count: len(threats)})
}

// GetThreatByID godoc
// @Summary Get threat by ID
// @Description Retrieve a specific threat by its ID
// @Tags threats
// @Produce json
// @Param threatId path string true "Threat ID"
// @Success 200 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Security ApiKeyAuth
// @Router /api/threats/{threatId} [get]
func (h *ThreatHandler) GetThreatByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	threat, err := h.service.GetThreatByID(vars["threatId"])
	if err != nil {
		logger.Warn("Threat not found: %s", vars["threatId"])
		utils.SendErrorResponse(w, "Threat not found", http.StatusNotFound, err.Error())
		return
	}
	utils.SendSuccessResponse(w, threat, nil)
}

// SearchThreats godoc
// @Summary Search threats
// @Description Search threats with filters (type, user)
// @Tags threats
// @Produce json
// @Param type query string false "Threat type"
// @Param user query string false "User ID"
// @Success 200 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security ApiKeyAuth
// @Router /api/threats/search [get]
func (h *ThreatHandler) SearchThreats(w http.ResponseWriter, r *http.Request) {
	filters := map[string]string{
		"type": r.URL.Query().Get("type"),
		"user": r.URL.Query().Get("user"),
	}
	
	threats, err := h.service.SearchThreats(filters)
	if err != nil {
		logger.Error("Threat search failed: %v", err)
		utils.SendErrorResponse(w, "Search failed", http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SendSuccessResponse(w, threats, &utils.Meta{Count: len(threats)})
}
