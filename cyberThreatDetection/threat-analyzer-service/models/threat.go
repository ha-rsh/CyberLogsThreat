package models

import "time"


type Log struct {
	ID					string					`json:"id,omitempty"`
	Timestamp			time.Time				`json:"timestamp"`
	UserID				string					`json:"userId"`
	IPAddress     		string    				`json:"ipAddress"`
	Action        		string    				`json:"action"`
	FileName      		*string   				`json:"fileName,omitempty"`
	DatabaseQuery 		*string   				`json:"databaseQuery,omitempty"`
}

type Threat struct {
	ID					string					`json:"id,omitempty"`
	Timestamp			time.Time				`json:"timestamp"`
	UserID				string					`json:"userId"`
	IPAddress			string					`json:"ipAddress"`
	Action				string					`json:"action"`
	FileName			*string					`json:"fileName,omitempty"`
	ThreatType			string					`json:"threatType"`
	Severity			string					`json:"severity"`
}

type AnalysisResult	struct {
	ThreatsDetected		int						`json:"threatsDetected"`
	Duration            string					`json:"duration"`
}