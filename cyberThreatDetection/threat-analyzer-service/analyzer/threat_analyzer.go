package analyzer

import (
	"cybersecuritySystem/shared/logger"
	"cybersecuritySystem/threat-analyzer-service/models"
	"sort"
)

type ThreatAnalyzer struct{}

func NewThreatAnalyzer() *ThreatAnalyzer {
	return &ThreatAnalyzer{}
}

func (a *ThreatAnalyzer) AnalyzeLogs(logs []models.Log) []models.Threat {
	logger.Info("Starting threat analysis on %d logs", len(logs))

	threats := []models.Threat{}
	sortedLogs := make([]models.Log, len(logs))
	copy(sortedLogs, logs)
	sort.Slice(sortedLogs, func(i, j int) bool {
		return sortedLogs[i].Timestamp.Before(sortedLogs[j].Timestamp)
	})

	userLogs := groupLogsByUser(sortedLogs)
	logger.Debug("Grouped logs by %d users", len(userLogs))

	for userID, userLogList := range userLogs {
		logger.Debug("Analyzing user: %s (%d logs)", userID, len(userLogList))

		threats = append(threats, DetectCredentialStuffing(userLogList)...)
		threats = append(threats, DetectPrivilegeEscalation(userLogList)...)
		threats = append(threats, DetectAccountTakeover(userLogList)...)
		threats = append(threats, DetectDataExfiltration(userLogList)...)
		threats = append(threats, DetectInsiderThreat(userLogList)...)

	}

	logger.Info("Threat analysis completed: %d threats detected", len(threats))
	return threats
}

func groupLogsByUser(logs []models.Log) map[string][]models.Log {
	userLogs := make(map[string][]models.Log)
	for _, log := range logs {
		userLogs[log.UserID] = append(userLogs[log.UserID], log)
	}
	return userLogs
}