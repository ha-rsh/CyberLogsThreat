package analyzer

import (
	"cybersecuritySystem/shared/constants"
	"cybersecuritySystem/shared/utils"
	"cybersecuritySystem/threat-analyzer-service/models"
	"sort"
	"time"
)

func DetectAccountTakeover(logs []models.Log) []models.Threat {
    threats := []models.Threat{}
    
    userLogs := make(map[string][]models.Log)
    for _, log := range logs {
        userLogs[log.UserID] = append(userLogs[log.UserID], log)
    }
    
    for _, userLogList := range userLogs {
        sort.Slice(userLogList, func(i, j int) bool {
            return userLogList[i].Timestamp.Before(userLogList[j].Timestamp)
        })
        
        loginEvents := []models.Log{}
        for _, log := range userLogList {
            if log.Action == constants.ActionLoginSuccess {
                loginEvents = append(loginEvents, log)
            }
        }
        
        for i := 0; i < len(loginEvents); i++ {
            for j := i + 1; j < len(loginEvents); j++ {
                login1 := loginEvents[i]
                login2 := loginEvents[j]
                
                timeDiff := login2.Timestamp.Sub(login1.Timestamp)
                
                if login1.IPAddress != login2.IPAddress && 
                   timeDiff <= 10*time.Minute {
                    
                  
                    
                    for _, fileLog := range userLogList {
                        if fileLog.Timestamp.After(login2.Timestamp) &&
                           fileLog.Action == constants.ActionFileAccess &&
                           fileLog.FileName != nil &&
                           utils.IsRestrictedFile(*fileLog.FileName) {
                            
                            threats = append(threats, models.Threat{
								ID:         utils.GenerateID(),
                                Timestamp:  fileLog.Timestamp,
                                UserID:     fileLog.UserID,
                                IPAddress:  fileLog.IPAddress,
                                Action:     fileLog.Action,
                                FileName:   fileLog.FileName,
                                ThreatType: constants.ThreatAccountTakeover,
                                Severity:   constants.SeverityCritical,
                            })
                            // Avoid duplicate threats
                            goto nextIPPair
                        }
                    }
                }
            }
            nextIPPair:
        }
    }
    
    return threats
}



// func DetectAccountTakeover(logs []models.Log) []models.Threat {
// 	threats := []models.Threat{}
// 	ipMap := make(map[string]time.Time)

// 	for _, log := range logs {
// 		if log.Action == constants.ActionLoginSuccess {
// 			for otherIP, otherTime := range ipMap {
// 				if otherIP != log.IPAddress {
// 					timeDiff := log.Timestamp.Sub(otherTime)
// 					if timeDiff <= 10*time.Minute && timeDiff >= 0 {
// 						for _, nextLog := range logs {
// 							if nextLog.Timestamp.After(log.Timestamp) &&
// 								nextLog.Action == constants.ActionFileAccess &&
// 								nextLog.FileName != nil &&
// 								utils.IsRestrictedFile(*nextLog.FileName) &&
// 								nextLog.UserID == log.UserID {
								
// 								threats = append(threats, models.Threat{
// 									Timestamp:  nextLog.Timestamp,
// 									UserID:     nextLog.UserID,
// 									IPAddress:  nextLog.IPAddress,
// 									Action:     nextLog.Action,
// 									FileName:   nextLog.FileName,
// 									ThreatType: constants.ThreatAccountTakeover,
// 									Severity:   constants.SeverityCritical,
// 								})
// 								break
// 							}
// 						}
// 					}
// 				}
// 			}
// 			ipMap[log.IPAddress] = log.Timestamp
// 		}
// 	}
// 	return threats
// }