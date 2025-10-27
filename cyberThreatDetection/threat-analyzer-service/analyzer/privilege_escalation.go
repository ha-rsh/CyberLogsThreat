package analyzer

import (
	"cybersecuritySystem/shared/constants"
	"cybersecuritySystem/shared/utils"
	"cybersecuritySystem/threat-analyzer-service/models"
	"time"
	"sort"
)


func DetectPrivilegeEscalation(logs []models.Log) []models.Threat {
    threats := []models.Threat{}
    userLogs := make(map[string][]models.Log)
    for _, log := range logs {
        userLogs[log.UserID] = append(userLogs[log.UserID], log)
    }
    
    for _, userLogList := range userLogs {
        sort.Slice(userLogList, func(i, j int) bool {
            return userLogList[i].Timestamp.Before(userLogList[j].Timestamp)
        })
        
        for i, log := range userLogList {
            if log.Action == constants.ActionLoginFailed {
                for j := i + 1; j < len(userLogList); j++ {
                    nextLog := userLogList[j]
                    timeDiff := nextLog.Timestamp.Sub(log.Timestamp)
                    if timeDiff > 5*time.Minute {
                        break
                    }
                    
                    if nextLog.DatabaseQuery != nil && 
                       utils.IsDangerousQuery(*nextLog.DatabaseQuery) {
                        
                        threats = append(threats, models.Threat{
                            ID:         utils.GenerateID(),
                            Timestamp:  nextLog.Timestamp,
                            UserID:     nextLog.UserID,
                            IPAddress:  nextLog.IPAddress,
                            Action:     nextLog.Action,
                            FileName:   nextLog.FileName,
                            ThreatType: constants.ThreatPrivilegeEscalation,
                            Severity:   constants.SeverityCritical,
                        })
                        break
                    }
                }
            }
        }
    }
    
    return threats
}



// func DetectPrivilegeEscalation(logs []models.Log) []models.Threat {
// 	threats := []models.Threat{}
// 	for i, log := range logs {
// 		if log.Action == constants.ActionLoginFailed {
// 			for j := i + 1; j < len(logs); j ++ {
// 				if logs[j].DatabaseQuery != nil && utils.IsDangerousQuery(*logs[j].DatabaseQuery) {
// 					timeDiff := logs[j].Timestamp.Sub(log.Timestamp)
// 					if timeDiff <= 5*time.Minute {
// 						threats = append(threats, models.Threat{
// 							Timestamp:  logs[j].Timestamp,
// 							UserID:     logs[j].UserID,
// 							IPAddress:  logs[j].IPAddress,
// 							Action:     logs[j].Action,
// 							FileName:   logs[j].FileName,
// 							ThreatType: constants.ThreatPrivilegeEscalation,
// 							Severity:   constants.SeverityCritical,
// 						})
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return threats
// }