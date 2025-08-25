package helper

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
)

type ElasticLog struct {
	Request  ElasticRequestLog  `json:"request"`
	Response ElasticResponseLog `json:"response"`
}

type ElasticRequestLog struct {
	RequestTime time.Time     `json:"request_time"`
	Method      string        `json:"method"`
	Uri         string        `json:"uri"`
	Proto       string        `json:"proto"`
	UserAgent   string        `json:"user_agent"`
	Referer     string        `json:"referer"`
	PostData    string        `json:"post_data"`
	ClientIP    string        `json:"client_ip"`
	LatencyTime time.Duration `json:"latency_time"`
	Level       string        `json:"level"`
	Key         string        `json:"key"`
}

type ElasticResponseLog struct {
	ResponseTime time.Time `json:"response_time"`
	StatusCode   int       `json:"status_code"`
	ResponseBody string    `json:"response_body"`
}

func SendLoggingToElastic(elastic *elasticsearch.Client, requestLog ElasticLog) {
	b, err := json.Marshal(requestLog)
	if err != nil {
		return
	}
	payload := string(b)

	u, err := uuid.NewRandom()
	if err != nil {
		return
	}

	// Instantiate a request object
	indexName := os.Getenv("ELASTIC_INDEX") + "-" + time.Now().Format("2006-01-02")
	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: u.String(),
		Body:       strings.NewReader(payload),
		Refresh:    "true",
	}

	// Return an API response object from request
	res, err := req.Do(context.Background(), elastic)
	if err != nil {
		log.Printf("IndexRequest ERROR: %s", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("%s ERROR indexing document ID=%s", res.Status(), u.String())
		return
	}
}
