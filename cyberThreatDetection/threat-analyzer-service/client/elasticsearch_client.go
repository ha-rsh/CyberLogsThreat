package client

import (
	"bytes"
	"context"
	"cybersecuritySystem/shared/constants"
	"cybersecuritySystem/shared/logger"
	"cybersecuritySystem/shared/utils"
	"cybersecuritySystem/threat-analyzer-service/models"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticsearchClient struct {
	client *elasticsearch.Client
}

func NewElasticsearchClient(esURL string) (*ElasticsearchClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Elasticsearch: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch returned error: %s", res.String())
	}

	client := &ElasticsearchClient{client: es}

	// Create threat index
	if err := client.createThreatIndex(); err != nil {
		return nil, fmt.Errorf("failed to create threat index: %v", err)
	}

	return client, nil
}

func (c *ElasticsearchClient) createThreatIndex() error {
	mapping := map[string]interface{} {
		"mappings": map[string]interface{} {
			"properties": map[string]interface{} {
				"timestamp":  map[string]interface{}{"type": "date"},
				"userId":     map[string]interface{}{"type": "keyword"},
				"ipAddress":  map[string]interface{}{"type": "ip"},
				"action":     map[string]interface{}{"type": "keyword"},
				"fileName":   map[string]interface{}{"type": "text"},
				"threatType": map[string]interface{}{"type": "keyword"},
				"severity":   map[string]interface{}{"type": "keyword"},
			},
		},
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(mapping)

	res, err := c.client.Indices.Create(
		constants.IndexThreats,
		c.client.Indices.Create.WithBody(&buf),
	)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		if !strings.Contains(string(bodyBytes), "resource_already_exists_exception") {
			return fmt.Errorf("error creating index: %s", res.String())
		}
	}

	return nil
}

func (c *ElasticsearchClient) GetAllLogs() ([]models.Log, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"size": 10000,
		"sort": []map[string]interface{}{
			{"timestamp": map[string]string{"order":"asc"}},
		},
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(query)

	res, err := c.client.Search(
		c.client.Search.WithContext(context.Background()),
		c.client.Search.WithIndex(constants.IndexLogs),
		c.client.Search.WithBody(&buf),
	)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error fetching logs: %s", res.String())
	}

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	logs := make([]models.Log, 0)

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		
		log := models.Log{}
		if id, ok := hitMap["_id"].(string); ok {
			log.ID = id
		}
		if ts, ok := source["timestamp"].(string); ok {
			log.Timestamp, _ = time.Parse(time.RFC3339, ts)
		}
		if userId, ok := source["userId"].(string); ok {
			log.UserID = userId
		}
		if ip, ok := source["ipAddress"].(string); ok {
			log.IPAddress = ip
		}
		if action, ok := source["action"].(string); ok {
			log.Action = action
		}
		if fileName, ok := source["fileName"].(string); ok && fileName != "" {
			log.FileName = &fileName
		}
		if dbQuery, ok := source["databaseQuery"].(string); ok && dbQuery != "" {
			log.DatabaseQuery = &dbQuery
		}

		logs = append(logs, log)
	}

	logger.Debug("Fetched %d logs from Elasticsearch", len(logs))
	return logs, nil
}


func (c *ElasticsearchClient) SaveThreat(threat *models.Threat) error {
	threat.ID = utils.GenerateID()
	
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(threat)

	req := esapi.IndexRequest{
		Index:      constants.IndexThreats,
		DocumentID: threat.ID,
		Body:       &buf,
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), c.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing threat: %s", res.String())
	}

	return nil
}

func (c *ElasticsearchClient) GetAllThreats() ([]map[string]interface{}, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"size": 10000,
		"sort": []map[string]interface{}{
			{"timestamp": map[string]string{"order": "desc"}},
		},
	}

	return c.search(query)
}

func (c *ElasticsearchClient) GetThreatByID(id string) (map[string]interface{}, error) {
	res, err := c.client.Get(constants.IndexThreats, id)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("threat not found")
	}

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	return result["_source"].(map[string]interface{}), nil
}

func (c *ElasticsearchClient) SearchThreats(query map[string]interface{}) ([]map[string]interface{}, error) {
	return c.search(query)
}

func (c *ElasticsearchClient) search(query map[string]interface{}) ([]map[string]interface{}, error) {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(query)

	res, err := c.client.Search(
		c.client.Search.WithContext(context.Background()),
		c.client.Search.WithIndex(constants.IndexThreats),
		c.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result map[string]interface{}
	json.NewDecoder(res.Body).Decode(&result)

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	documents := make([]map[string]interface{}, 0)

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})
		source["id"] = hitMap["_id"].(string)
		documents = append(documents, source)
	}

	return documents, nil
}