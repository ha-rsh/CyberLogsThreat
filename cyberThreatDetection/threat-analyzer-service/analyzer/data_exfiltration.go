package analyzer

import (
	"cybersecuritySystem/shared/constants"
	"cybersecuritySystem/shared/utils"
	"cybersecuritySystem/threat-analyzer-service/models"
	"time"
	"sort"
)


func DetectDataExfiltration(logs []models.Log) []models.Threat {
    threats := []models.Threat{}
    
    userLogs := make(map[string][]models.Log)
    for _, log := range logs {
        if log.Action == constants.ActionFileAccess && 
           log.FileName != nil && 
           utils.IsRestrictedFile(*log.FileName) {
            userLogs[log.UserID] = append(userLogs[log.UserID], log)
        }
    }
    
    for _, restrictedAccesses := range userLogs {
        sort.Slice(restrictedAccesses, func(i, j int) bool {
            return restrictedAccesses[i].Timestamp.Before(restrictedAccesses[j].Timestamp)
        })
        
        for i := 0; i < len(restrictedAccesses); i++ {
            baseTime := restrictedAccesses[i].Timestamp
            accessCount := 1
            lastAccess := restrictedAccesses[i]
            
            for j := i + 1; j < len(restrictedAccesses); j++ {
                timeDiff := restrictedAccesses[j].Timestamp.Sub(baseTime)
                
                if timeDiff <= 30*time.Second {
                    accessCount++
                    lastAccess = restrictedAccesses[j]
                } else {
                    break
                }
            }
            
            if accessCount >= 2 {
                threats = append(threats, models.Threat{
					ID:         utils.GenerateID(),
                    Timestamp:  lastAccess.Timestamp,
                    UserID:     lastAccess.UserID,
                    IPAddress:  lastAccess.IPAddress,
                    Action:     lastAccess.Action,
                    FileName:   lastAccess.FileName,
                    ThreatType: constants.ThreatDataExfiltration,
                    Severity:   constants.SeverityHigh,
                })
                
                i += accessCount - 1
                break
            }
        }
    }
    
    return threats
}




// func DetectDataExfiltration(logs []models.Log) []models.Threat {
// 	threats := []models.Threat{}
// 	fileAccesses := []models.Log{}

// 	for _, log := range logs {
// 		if log.Action == constants.ActionFileAccess && 
// 			log.FileName != nil && 
// 			utils.IsRestrictedFile(*log.FileName) {
// 			fileAccesses = append(fileAccesses, log)

// 			if len(fileAccesses) >= 2 {
// 				lastAccess := fileAccesses[len(fileAccesses)-1]

// 				for idx := len(fileAccesses) - 2; idx >= 0; idx-- {
// 					if lastAccess.Timestamp.Sub(fileAccesses[idx].Timestamp) > 30*time.Second {
// 						fileAccesses = fileAccesses[idx+1:]
// 						break
// 					}
// 				}

// 				if len(fileAccesses) >= 3 {
// 					threats = append(threats, models.Threat{
// 						Timestamp:  lastAccess.Timestamp,
// 						UserID:     lastAccess.UserID,
// 						IPAddress:  lastAccess.IPAddress,
// 						Action:     lastAccess.Action,
// 						FileName:   lastAccess.FileName,
// 						ThreatType: constants.ThreatDataExfiltration,
// 						Severity:   constants.SeverityHigh,
// 					})
// 					fileAccesses = []models.Log{}
// 				}
// 			}
// 		} else {
// 			fileAccesses = []models.Log{}
// 		}
// 	}
// 	return threats
// }