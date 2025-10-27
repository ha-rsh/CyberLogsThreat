package service

import (
	"cybersecuritySystem/log-ingester-service/models"
	"cybersecuritySystem/log-ingester-service/repository"
	"cybersecuritySystem/shared/utils"
	"time"
)


type LogService struct {
	repo *repository.LogRepository
}

func NewLogService(repo *repository.LogRepository) *LogService {
	return &LogService{repo: repo}
}

func (s *LogService) CreateLog(log *models.Log) error {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	if err := utils.ValidateLog(log.UserID, log.IPAddress, log.Action); err != nil {
		return err
	}

	log.ID = utils.GenerateID()
	return s.repo.Create(log)
}

func (s *LogService) GetAllLogs() ([]map[string]interface{}, error) {
	return s.repo.FindAll()
}

func (s *LogService) GetLogByID(id string) (map[string]interface{}, error) {
	return s.repo.FindByID(id)
}

func (s *LogService) SearchLogs(filters map[string]string) ([]map[string]interface{}, error) {
	query := buildSearchQuery(filters)
	return s.repo.Search(query)
}

func (s *LogService) BulkCreateLogs(logs []models.Log) error {
	for i := range logs {
		if logs[i].Timestamp.IsZero() {
			logs[i].Timestamp = time.Now()
		}
		logs[i].ID = utils.GenerateID()
	}
	
	return s.repo.BulkCreate(logs)
}

func buildSearchQuery(filters map[string]string) map[string]interface{} {
	must := []map[string]interface{}{}
	
	if userId, ok := filters["userId"]; ok && userId != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{"userId": userId},
		})
	}
	
	if action, ok := filters["action"]; ok && action != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{"action": action},
		})
	}

	if ipAddress, ok := filters["ipAddress"]; ok && ipAddress != "" {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{"ipAddress": ipAddress},
		})
	}
	
	if startTime, ok := filters["startTime"]; ok && startTime != "" {
		if endTime, ok2 := filters["endTime"]; ok2 && endTime != "" {
			must = append(must, map[string]interface{}{
				"range": map[string]interface{}{
					"timestamp": map[string]interface{}{
						"gte": startTime,
						"lte": endTime,
					},
				},
			})
		}
	}
	
	if len(must) == 0 {
		return map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
			"size": 1000,
			"sort": []map[string]interface{}{{"timestamp": map[string]string{"order": "desc"}}},
		}
	}
	
	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"size": 1000,
		"sort": []map[string]interface{}{{"timestamp": map[string]string{"order": "desc"}}},
	}
}