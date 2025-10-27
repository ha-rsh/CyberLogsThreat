package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

// APIResponse represents the standard API response format
// @Description Standard API response wrapper
type APIResponse struct {
	Success 		bool        	`json:"success"`
	Data    		interface{} 	`json:"data,omitempty"`
	Error   		*APIError   	`json:"error,omitempty"`
	Meta    		*Meta       	`json:"meta,omitempty"`
}

// APIError represents an error in the API response
// @Description API error details
type APIError struct {
	Code    		int    			`json:"code"`
	Message 		string 			`json:"message"`
	Details 		string 			`json:"details,omitempty"`
}

// Meta represents metadata in the API response
// @Description Response metadata with pagination and timing info
type Meta struct {
	Page      		int   			`json:"page,omitempty"`
	Limit     		int   			`json:"limit,omitempty"`
	Total     		int   			`json:"total,omitempty"`
	Count     		int   			`json:"count,omitempty"`
	Timestamp 		int64 			`json:"timestamp"`
}

func SendSuccessResponse(w http.ResponseWriter, data interface{}, meta *Meta) {
	w.Header().Set("Content-Type", "application/json")
	
	if meta == nil {
		meta = &Meta{}
	}
	meta.Timestamp = time.Now().Unix()
	
	response := APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
	
	json.NewEncoder(w).Encode(response)
}

func SendErrorResponse(w http.ResponseWriter, message string, code int, details ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	apiError := &APIError{
		Code:    code,
		Message: message,
	}
	
	if len(details) > 0 {
		apiError.Details = details[0]
	}
	
	response := APIResponse{
		Success: false,
		Error:   apiError,
		Meta: &Meta{
			Timestamp: time.Now().Unix(),
		},
	}
	
	json.NewEncoder(w).Encode(response)
}