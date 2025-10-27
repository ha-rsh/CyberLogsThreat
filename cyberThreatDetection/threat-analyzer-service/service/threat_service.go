package service

import (
	"cybersecuritySystem/shared/logger"
	"cybersecuritySystem/threat-analyzer-service/analyzer"
	"cybersecuritySystem/threat-analyzer-service/client"
	"cybersecuritySystem/threat-analyzer-service/models"
	"time"
)

type ThreatService struct {
	client   *client.ElasticsearchClient
	analyzer *analyzer.ThreatAnalyzer
}

func NewThreatService(client *client.ElasticsearchClient) *ThreatService {
	return &ThreatService{
		client:   client,
		analyzer: analyzer.NewThreatAnalyzer(),
	}
}

func (s *ThreatService) AnalyzeThreats() (models.AnalysisResult, error) {
	start := time.Now()
	logger.Info("Starting threat analysis")
	logs, err := s.client.GetAllLogs()
	if err != nil {
		logger.Error("Failed to fetch logs: %v", err)
		return models.AnalysisResult{}, err
	}

	if len(logs) == 0 {
		logger.Warn("NO logs foundfor analysis")
		return models.AnalysisResult{ThreatsDetected: 0, Duration: time.Since(start).String()}, nil
	}

	threats := s.analyzer.AnalyzeLogs(logs)
	for i := range threats {
		if err := s.client.SaveThreat(&threats[i]); err != nil {
			logger.Error("Failed to save threat: %v", err)
		}
	}

	duration := time.Since(start)
	logger.Success("Analysis completed: %d threats detected in %s", len(threats), duration)

	return models.AnalysisResult{
		ThreatsDetected: len(threats),
		Duration: duration.String(),
	}, nil
}

func (s *ThreatService) GetAllThreats() ([]map[string]interface{}, error) {
	return s.client.GetAllThreats()
}

func (s *ThreatService) GetThreatByID(id string) (map[string]interface{}, error) {
	return s.client.GetThreatByID(id)
}

func (s *ThreatService) SearchThreats(filters map[string]string) ([]map[string]interface{}, error) {
	query := buildSearchQuery(filters)
	return s.client.SearchThreats(query)
}

func buildSearchQuery(filters map[string]string) map[string]interface{} {
	must := []map[string]interface{}{}
	if threatType, ok := filters["type"]; ok && threatType != "" {
		must = append(must, map[string]interface{}{"match": map[string]interface{}{"threatType": threatType}})
	}

	if userId, ok := filters["user"]; ok && userId != "" {
		must = append(must, map[string]interface{}{"match": map[string]interface{}{"userId": userId}})
	}

	if len(must) == 0 {
		return map[string]interface{}{"query": map[string]interface{}{"match_all": map[string]interface{}{}}, "size": 1000}
	}
	return map[string]interface{}{"query": map[string]interface{}{"bool": map[string]interface{}{"must": must}}, "size": 1000}

}