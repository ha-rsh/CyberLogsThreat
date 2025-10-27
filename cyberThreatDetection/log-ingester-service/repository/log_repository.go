package repository

import (
	"bytes"
	"context"
	"cybersecuritySystem/log-ingester-service/models"
	"cybersecuritySystem/shared/constants"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type LogRepository struct {
	client *elasticsearch.Client
}


func NewLogRepository(esURL string) (*LogRepository, error) {
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

	repo := &LogRepository{client: es}

	if err := repo.createIndex(); err != nil {
		return nil, fmt.Errorf("failed to create index: %v", err)
	}

	return repo, nil

}

func (r *LogRepository) createIndex() error {
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"timestamp":     map[string]interface{}{"type": "date"},
				"userId":        map[string]interface{}{"type": "keyword"},
				"ipAddress":     map[string]interface{}{"type": "ip"},
				"action":        map[string]interface{}{"type": "keyword"},
				"fileName":      map[string]interface{}{"type": "text"},
				"databaseQuery": map[string]interface{}{"type": "text"},
			},
		},
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(mapping)

	res, err := r.client.Indices.Create(
		constants.IndexLogs,
		r.client.Indices.Create.WithBody(&buf),
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

func (r *LogRepository) Create(log *models.Log) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(log); err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      constants.IndexLogs,
		DocumentID: log.ID,
		Body:       &buf,
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

func (r *LogRepository) FindAll() ([]map[string]interface{}, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"size": 10000,
	}

	return r.search(query)
}

func (r *LogRepository) FindByID(id string) (map[string]interface{}, error) {
	res, err := r.client.Get(constants.IndexLogs, id)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("document not found")
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result["_source"].(map[string]interface{}), nil
}

func (r *LogRepository) Search(query map[string]interface{}) ([]map[string]interface{}, error) {
	return r.search(query)
}

func (r *LogRepository) BulkCreate(logs []models.Log) error {
	var buf bytes.Buffer

	for _, log := range logs {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": constants.IndexLogs,
				"_id":    log.ID,
			},
		}
		json.NewEncoder(&buf).Encode(meta)
		json.NewEncoder(&buf).Encode(log)
	}

	res, err := r.client.Bulk(bytes.NewReader(buf.Bytes()), r.client.Bulk.WithIndex(constants.IndexLogs))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error in bulk indexing: %s", res.String())
	}

	return nil
}

func (r *LogRepository) search(query map[string]interface{}) ([]map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(context.Background()),
		r.client.Search.WithIndex(constants.IndexLogs),
		r.client.Search.WithBody(&buf),
		r.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

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
