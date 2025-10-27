package analyzer

import (
	"cybersecuritySystem/shared/constants"
	"cybersecuritySystem/shared/utils"
	"cybersecuritySystem/threat-analyzer-service/models"
)

func DetectInsiderThreat(logs []models.Log) []models.Threat {
    threats := []models.Threat{}
    
    for _, log := range logs {
        if log.Action == constants.ActionFileAccess {
            hour := log.Timestamp.Hour()
            
            if hour >= 2 && hour < 5 {
                threats = append(threats, models.Threat{
					ID:         utils.GenerateID(),
                    Timestamp:  log.Timestamp,
                    UserID:     log.UserID,
                    IPAddress:  log.IPAddress,
                    Action:     log.Action,
                    FileName:   log.FileName,
                    ThreatType: constants.ThreatInsiderThreat,
                    Severity:   constants.SeverityMedium,
                })
            }
        }
    }
    
    return threats
}