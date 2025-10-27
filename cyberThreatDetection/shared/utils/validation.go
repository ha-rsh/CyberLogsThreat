package utils

import (
	// "log"
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
)

func ValidateLog(userId, ipAddress, action string) error {
	if userId == "" {
		return NewValidationError("userId is required")
	}
	
	if ipAddress == "" {
		return NewValidationError("ipAddress is required")
	}
	
	if !IsValidIP(ipAddress) {
		return NewValidationError("invalid IP address format")
	}
	
	if action == "" {
		return NewValidationError("action is required")
	}
	
	return nil
}

func IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsRestrictedFile(filename string) bool {
    if filename == "" {
        return false
    }
    
    filename = strings.ToLower(filename)
    
    // Based on your Elasticsearch data, check for these patterns:
    restrictedPatterns := []string{
        "payroll.csv",
        "/secure/payroll.csv", 
        "secure/payroll.csv",
        "design.pdf",
        "/confidential/design.pdf",
        "confidential/design.pdf", 
        "db_dump.sql",
        "/db_dump.sql",
        "system.log",
        "/logs/system.log",
        "logs/system.log",
    }
    
    for _, pattern := range restrictedPatterns {
        if strings.Contains(filename, pattern) {
            fmt.Printf("[IsRestrictedFile] ✅ File '%s' matches pattern '%s'", filename, pattern)
            return true
        }
    }
    
    fmt.Printf("[IsRestrictedFile] ❌ File '%s' is not restricted", filename)
    return false
}


func IsDangerousQuery(query string) bool {
    if query == "" {
        return false
    }
    
    query = strings.ToLower(query)
    dangerousPatterns := []string{
        "drop table",
        "delete from",
        "insert into admins",
        "insert into users",
        "update users set",
        "grant all privileges",
        "create user",
    }
    
    for _, pattern := range dangerousPatterns {
        if strings.Contains(query, pattern) {
            return true
        }
    }
    return false
}


// func IsRestrictedFile(fileName string) bool {
// 	if fileName == "" {
// 		return false
// 	}
	
// 	restrictedPaths := []string{"/secure/", "/confidential/", "/admin/", "payroll", "financial", "secret"}
// 	lowerFileName := strings.ToLower(fileName)
	
// 	for _, path := range restrictedPaths {
// 		if strings.Contains(lowerFileName, path) {
// 			return true
// 		}
// 	}

// 	log.Printf("[IsRestrictedFile] File '%s' is NOT restricted", fileName)
	
// 	return false
// }

// func IsDangerousQuery(query string) bool {
// 	if query == "" {
// 		return false
// 	}
	
// 	dangerousKeywords := []string{"DROP", "DELETE", "INSERT INTO admins", "UPDATE admins", "GRANT"}
// 	upperQuery := strings.ToUpper(query)
	
// 	for _, keyword := range dangerousKeywords {
// 		if strings.Contains(upperQuery, keyword) {
// 			return true
// 		}
// 	}
	
// 	return false
// }

func GenerateID() string {
	return uuid.New().String()
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}