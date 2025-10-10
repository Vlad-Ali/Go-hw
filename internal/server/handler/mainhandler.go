package handler

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	version = "v1.0.0"
)

type MainHandler struct {
}

func NewMainHandler() *MainHandler {
	return &MainHandler{}
}

func (t *MainHandler) GetVersion(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("Got request for %s", r.URL.Path)
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte(version))
	if err != nil {
		log.Printf("Error writing response: %v", err)
		return
	}
	log.Printf("Successfully wrote response")
}

func (t *MainHandler) DecodeString(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("Got request for %s", r.URL.Path)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	var data struct {
		InputString string `json:"inputString"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error parsing body: %v", err)
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(data.InputString)
	if err != nil {
		log.Printf("Error decoding body: %v", err)
		http.Error(w, "Invalid inputString", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"outputString": string(decoded),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
	log.Println("Successfully encoded response")
}

func (t *MainHandler) GetHardOp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("Got request for %s", r.URL.Path)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	seconds := rnd.Intn(11) + 10
	time.Sleep(time.Duration(seconds) * time.Second)
	randStatus := rnd.Intn(2) + 1
	if randStatus == 1 {
		log.Printf("Successfully executed hard operation")
		w.WriteHeader(http.StatusOK)
	} else {
		log.Printf("Failed to execute hard operation")
		w.WriteHeader(rnd.Intn(12) + 500)
	}
}
