package models

import "time"


type Log struct {
	ID					string			`json:"id,omitempty"`
	Timestamp			time.Time		`json:"timestamp"`
	UserID				string			`json:"userId" validate:"required"`
	IPAddress			string			`json:"ipAddress" validate:"required"`
	Action      		string			`json:"action" validate:"required,oneof=loginSuccess loginFailed fileAccess databaseQuery"`
	FileName    		*string  		`json:"fileName,omitempty"`
	DatabaseQuery 		*string 		`json:"databaseQuery,omitempty"`
}

type FileUploadRequest struct {
	File				[]byte			`json:"file"`	
}

type FileUploadResponse struct {
	FileName			string			`json:"fileName"`
	LogsUploaded		int 			`json:"logsUploaded"`
}