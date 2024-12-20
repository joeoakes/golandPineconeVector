package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	apiKey   = "YOUR_PINECONE_API_KEY"              // Replace with your Pinecone API key
	indexURL = "https://your-index-url/pinecone.io" // Replace with your Pinecone index URL
)

type Vector struct {
	ID       string                 `json:"id"`
	Values   []float64              `json:"values"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type UpsertRequest struct {
	Vectors []Vector `json:"vectors"`
}

type QueryRequest struct {
	TopK            int       `json:"topK"`
	Values          []float64 `json:"values,omitempty"`
	ID              string    `json:"id,omitempty"`
	IncludeMetadata bool      `json:"includeMetadata"`
}

type QueryResponse struct {
	Matches []struct {
		ID       string                 `json:"id"`
		Score    float64                `json:"score"`
		Metadata map[string]interface{} `json:"metadata"`
	} `json:"matches"`
}

func main() {
	// Example vector
	vector := Vector{
		ID:     "example-vector-1",
		Values: []float64{0.1, 0.2, 0.3},
		Metadata: map[string]interface{}{
			"name": "example",
			"type": "demo",
		},
	}

	// 1. Create (Upsert) a vector
	upsertVector(vector)

	// 2. Read (Query) vectors
	queryVector([]float64{0.1, 0.2, 0.3}, 3)

	// 3. Update (Re-Upsert with modified data)
	vector.Metadata["type"] = "updated-demo"
	upsertVector(vector)

	// 4. Delete a vector
	deleteVector("example-vector-1")
}

func upsertVector(vector Vector) {
	url := fmt.Sprintf("%s/vectors/upsert", indexURL)

	requestBody, err := json.Marshal(UpsertRequest{
		Vectors: []Vector{vector},
	})
	if err != nil {
		log.Fatalf("Error marshaling upsert request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Error creating upsert request: %v", err)
	}
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error executing upsert request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Upsert response: %s\n", body)
}

func queryVector(values []float64, topK int) {
	url := fmt.Sprintf("%s/query", indexURL)

	requestBody, err := json.Marshal(QueryRequest{
		TopK:            topK,
		Values:          values,
		IncludeMetadata: true,
	})
	if err != nil {
		log.Fatalf("Error marshaling query request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Error creating query request: %v", err)
	}
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error executing query request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var response QueryResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Fatalf("Error unmarshaling query response: %v", err)
	}

	fmt.Printf("Query response: %+v\n", response)
}

func deleteVector(id string) {
	url := fmt.Sprintf("%s/vectors/delete", indexURL)

	requestBody, err := json.Marshal(map[string]interface{}{
		"ids": []string{id},
	})
	if err != nil {
		log.Fatalf("Error marshaling delete request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Error creating delete request: %v", err)
	}
	req.Header.Set("Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error executing delete request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Delete response: %s\n", body)
}
