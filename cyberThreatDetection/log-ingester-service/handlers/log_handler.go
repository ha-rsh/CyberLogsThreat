package handlers

import (
	"cybersecuritySystem/log-ingester-service/models"
	"cybersecuritySystem/log-ingester-service/service"
	"cybersecuritySystem/shared/logger"
	"cybersecuritySystem/shared/utils"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type LogHandler struct {
	service *service.LogService
}

func NewLogHandler(service *service.LogService) *LogHandler {
	return &LogHandler{service: service}
}

// UploadCSV godoc
// @Summary Upload CSV log file
// @Description Upload a CSV file containing logs and import them into the system
// @Tags logs
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file to upload"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security ApiKeyAuth
// @Router /api/logs/upload [post]
func (h *LogHandler) UploadCSV(w http.ResponseWriter, r *http.Request) {
    logger.Info("Starting CSV upload...")
    
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        utils.SendErrorResponse(w, "Failed to parse form", http.StatusBadRequest, err.Error())
        return
    }

    file, handler, err := r.FormFile("file")
    if err != nil {
        logger.Error("Error reading file: %v", err)
        utils.SendErrorResponse(w, "Failed to read file", http.StatusBadRequest, err.Error())
        return
    }
    defer file.Close()

    // Validate file extension
    if !strings.HasSuffix(strings.ToLower(handler.Filename), ".csv") {
        utils.SendErrorResponse(w, "Invalid file format. Only CSV files are allowed", http.StatusBadRequest, "")
        return
    }

    reader := csv.NewReader(file)
    reader.FieldsPerRecord = -1
    reader.TrimLeadingSpace = true

    // Read headers
    headers, err := reader.Read()
    if err != nil {
        logger.Error("Error reading CSV headers: %v", err)
        utils.SendErrorResponse(w, "Invalid CSV format", http.StatusBadRequest, err.Error())
        return
    }

    if len(headers) == 0 {
        utils.SendErrorResponse(w, "Empty CSV file", http.StatusBadRequest, "")
        return
    }

    var (
        logs           = make([]models.Log, 0, 100)
        totalProcessed = 0
        totalErrors    = 0
        batchSize      = 100
    )

    // Process CSV records
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        
        if err != nil {
            totalErrors++
            continue
        }

        parsedLog := parseCSVRecord(headers, record)
        
        // Validate required fields
        if parsedLog.UserID == "" || parsedLog.IPAddress == "" || parsedLog.Action == "" {
            totalErrors++
            continue
        }
        
        logs = append(logs, parsedLog)

        // Process batch
        if len(logs) >= batchSize {
            batchProcessed, batchErr := h.processBatch(logs, 0)
            totalProcessed += batchProcessed
            
            if batchErr != nil {
                totalErrors++
            }
            
            logs = logs[:0]
        }
    }

    // Process final batch
    if len(logs) > 0 {
        batchProcessed, batchErr := h.processBatch(logs, 0)
        totalProcessed += batchProcessed
        
        if batchErr != nil {
            totalErrors++
        }
    }

    logger.Info("Upload complete: %s, processed: %d, errors: %d", 
        handler.Filename, totalProcessed, totalErrors)

    // Handle many errors
    if totalProcessed == 0 && totalErrors > 0 {
        utils.SendErrorResponse(w, "Upload failed - no valid records processed", http.StatusBadRequest, "")
        return
    }

    response := models.FileUploadResponse{
        FileName:     handler.Filename,
        LogsUploaded: totalProcessed,
    }

    utils.SendSuccessResponse(w, response, &utils.Meta{Count: totalProcessed})
}

func (h *LogHandler) processBatch(logs []models.Log, batchNum int) (int, error) {
    if len(logs) == 0 {
        return 0, nil
    }

    startTime := time.Now()
    err := h.service.BulkCreateLogs(logs)
    duration := time.Since(startTime)
    
    logger.Info("Batch %d processing took %v", batchNum, duration)
    
    if err != nil {
        return 0, fmt.Errorf("failed to process batch of %d logs: %w", len(logs), err)
    }
    
    return len(logs), nil
}


// CreateLog godoc
// @Summary Create a new log
// @Description Create a new log entry in the system
// @Tags logs
// @Accept json
// @Produce json
// @Param log body models.Log true "Log object to create"
// @Success 201 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security ApiKeyAuth
// @Router /api/logs [post]
func (h *LogHandler) CreateLog(w http.ResponseWriter, r *http.Request) {
	var log models.Log
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.CreateLog(&log); err != nil {
		if _, ok := err.(*utils.ValidationError); ok {
			utils.SendErrorResponse(w, "Validation error", http.StatusBadRequest, err.Error())
			return
		}
		utils.SendErrorResponse(w, "Failed to create log", http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	utils.SendSuccessResponse(w, log, nil)
}

// GetLogs godoc
// @Summary Get all logs
// @Description Retrieve all logs from the system
// @Tags logs
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security ApiKeyAuth
// @Router /api/logs [get]
func (h *LogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.service.GetAllLogs()
	if err != nil {
		utils.SendErrorResponse(w, "Failed to retrieve logs", http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, logs, &utils.Meta{Count: len(logs)})
}

// GetLogByID godoc
// @Summary Get log by ID
// @Description Retrieve a specific log by its ID
// @Tags logs
// @Produce json
// @Param logId path string true "Log ID"
// @Success 200 {object} utils.APIResponse
// @Failure 404 {object} utils.APIResponse
// @Security ApiKeyAuth
// @Router /api/logs/{logId} [get]
func (h *LogHandler) GetLogByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logID := vars["logId"]

	log, err := h.service.GetLogByID(logID)
	if err != nil {
		utils.SendErrorResponse(w, "Log not found", http.StatusNotFound, err.Error())
		return
	}

	utils.SendSuccessResponse(w, log, nil)
}

// SearchLogs godoc
// @Summary Search logs
// @Description Search logs with filters (userId, action, time range)
// @Tags logs
// @Produce json
// @Param userId query string false "User ID"
// @Param action query string false "Action type"
// @Param startTime query string false "Start time (ISO 8601)"
// @Param endTime query string false "End time (ISO 8601)"
// @Success 200 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Security ApiKeyAuth
// @Router /api/logs/search [get]
func (h *LogHandler) SearchLogs(w http.ResponseWriter, r *http.Request) {
	filters := map[string]string{
		"userId":    r.URL.Query().Get("userId"),
		"action":    r.URL.Query().Get("action"),
		"startTime": r.URL.Query().Get("startTime"),
		"endTime":   r.URL.Query().Get("endTime"),
	}

	logs, err := h.service.SearchLogs(filters)
	if err != nil {
		utils.SendErrorResponse(w, "Search failed", http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, logs, &utils.Meta{Count: len(logs)})
}

func parseCSVRecord(headers, record []string) models.Log {
	logEntry := models.Log{}

	for i, header := range headers {
		if i >= len(record) {
			break
		}

		value := record[i]
		if value == "" || value == "NaN" {
			continue
		}

		headerLower := strings.ToLower(strings.TrimSpace(header))

		switch headerLower {
		case "timestamp":
			formats := []string{
				time.RFC3339,
				"2006-01-02 15:04:05",
				"1/2/2006 15:04",
				"2006-01-02T15:04:05Z", 
			}
			
			parsed := false
			for _, format := range formats {
				if t, err := time.Parse(format, value); err == nil {
					logEntry.Timestamp = t
					parsed = true
					break
				}
			}
			
			if !parsed {
				logger.Warn("[parseCSVRecord] Failed to parse timestamp: '%s'\n", value)
			}
			
		case "user_id", "userid":
			logEntry.UserID = value
			
		case "ip_address", "ipaddress":
			logEntry.IPAddress = value
			
		case "action":
			logEntry.Action = normalizeAction(value)
			
		case "file_name", "filename":
			if value != "" {
				logEntry.FileName = &value
			}
			
		case "database_query", "databasequery":
			if value != "" {
				logEntry.DatabaseQuery = &value
			}
		}
	}

	if logEntry.Timestamp.IsZero() {
		logEntry.Timestamp = time.Now()
	}

	return logEntry
}

func normalizeAction(action string) string {
	action = strings.ToLower(strings.TrimSpace(action))
	parts := strings.Split(action, "_")

	if len(parts) == 1 {
		return action
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result += strings.ToUpper(string(parts[i][0])) + parts[i][1:]
		}
	}

	return result
}