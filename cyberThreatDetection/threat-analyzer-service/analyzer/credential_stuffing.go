package analyzer

import (
	"cybersecuritySystem/shared/constants"
	"cybersecuritySystem/shared/utils"
	"cybersecuritySystem/threat-analyzer-service/models"
	"fmt"
	"sort"
)

func DetectCredentialStuffing(logs []models.Log) []models.Threat {
    threats := []models.Threat{}
    
    userLogs := make(map[string][]models.Log)
    for _, log := range logs {
        userLogs[log.UserID] = append(userLogs[log.UserID], log)
    }
    
    for userID, userLogList := range userLogs {
        sort.Slice(userLogList, func(i, j int) bool {
            return userLogList[i].Timestamp.Before(userLogList[j].Timestamp)
        })
        
        failedLogins := []models.Log{}
        
        for i, log := range userLogList {
            if log.Action == constants.ActionLoginFailed {
                failedLogins = append(failedLogins, log)
                
            } else if log.Action == constants.ActionLoginSuccess {
                if len(failedLogins) >= 2 {
                    fmt.Printf("User %s: %d login success after logins failed", 
                        userID, len(failedLogins))
                    
                    fmt.Printf("Access sensitive file after login for user %s", userID)
                    
                    foundSensitiveAccess := false
                    for j := i + 1; j < len(userLogList); j++ {
                        fileLog := userLogList[j]
                        
                        if fileLog.Action == constants.ActionFileAccess {
                            if fileLog.FileName != nil {
                                
                                if utils.IsRestrictedFile(*fileLog.FileName) {
                                    foundSensitiveAccess = true
                                    threats = append(threats, models.Threat{
                                        ID:         utils.GenerateID(),
                                        Timestamp:  fileLog.Timestamp,
                                        UserID:     fileLog.UserID,
                                        IPAddress:  fileLog.IPAddress,
                                        Action:     fileLog.Action,
                                        FileName:   fileLog.FileName,
                                        ThreatType: constants.ThreatCredentialStuffing,
                                        Severity:   constants.SeverityHigh,
                                    })
                                    break
                                } else {
                                    fmt.Printf("File %s is NOT restricted", *fileLog.FileName)
                                }
                            }
                        }
                    }
                    
                    if !foundSensitiveAccess {
                        fmt.Printf("No sensitive file access found for user %s after successful login", userID)
                    }
                }
                
                failedLogins = []models.Log{}
            }
        }
    }
    
    return threats
}


// func DetectCredentialStuffing(logs []models.Log) []models.Threat {
// 	threats := []models.Threat{}
// 	failedLogins := []models.Log{}

// 	for i, log := range logs {
// 		if log.Action == constants.ActionLoginFailed {
// 			failedLogins = append(failedLogins, log)
// 		} else if log.Action == constants.ActionLoginSuccess && len(failedLogins) > 1 {
// 			for j := i + 1; j < len(logs); j++ {
// 				if logs[j].Action == constants.ActionFileAccess && logs[j].FileName != nil && utils.IsRestrictedFile(*logs[j].FileName) {
// 					threats = append(threats, models.Threat{
// 						Timestamp: logs[j].Timestamp,
// 						UserID:     logs[j].UserID,
// 						IPAddress:  logs[j].IPAddress,
// 						Action:     logs[j].Action,
// 						FileName:   logs[j].FileName,
// 						ThreatType: constants.ThreatCredentialStuffing,
// 						Severity:   constants.SeverityHigh,
// 					})
// 					break
// 				}
// 			}
// 			failedLogins = []models.Log{}
// 		}
// 	}
// 	return  threats
// }