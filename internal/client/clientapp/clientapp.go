package clientapp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	ErrCreatingRequest = errors.New("creating request failed")
	ErrRequestFailed   = errors.New("request failed")
	ErrReadingResponse = errors.New("reading response failed")
)

type ClientApp struct {
	client  *http.Client
	baseURL string
}

func NewClientApp(baseURL string) *ClientApp {
	return &ClientApp{client: &http.Client{
		Timeout: 30 * time.Second,
	}, baseURL: baseURL}
}

func (c *ClientApp) GetVersionEndpoint() error {
	log.Println("GetVersionEndpoint called")

	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/version", nil)
	if err != nil {
		log.Printf("Creating request failed: %v", err)
		return ErrCreatingRequest
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return ErrRequestFailed
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading response body failed: %v", err)
		return ErrReadingResponse
	}
	log.Println(string(body))
	return nil
}

func (c *ClientApp) PostDecodeEndpoint(inputString string) error {
	log.Println("PostDecodeEndpoint called")

	requestBody := map[string]string{"inputString": inputString}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Marshalling request body failed: %v", err)
		return errors.New("Marshalling request body failed.")
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/decode", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Creating request failed: %v", err)
		return ErrCreatingRequest
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return ErrRequestFailed
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading response body failed: %v", err)
		return ErrReadingResponse
	}

	var data struct {
		OutputString string `json:"outputString"`
	}
	if err = json.Unmarshal(body, &data); err != nil {
		log.Printf("Error parsing body failed: %v", err)
		return errors.New("Error parsing body")
	}
	log.Println(data.OutputString)
	return nil
}

func (c *ClientApp) GetHardOpEndpoint() error {
	log.Printf("GetHardOpEndpoint called")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/hard-op", nil)
	if err != nil {
		log.Printf("Creating request failed: %v", err)
		return ErrCreatingRequest
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		log.Printf("Request failed: %v", err)
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Request timeout exceeded")
		}
		return ErrRequestFailed
	}
	defer resp.Body.Close()
	log.Printf("Request timeout not exceeded %v", resp.StatusCode)
	return nil
}
