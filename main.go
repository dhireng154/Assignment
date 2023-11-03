package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi"
)

// The Server struct is the primary server, featuring an in-memory data store
type Server struct {
	DataStore *InMemoryDataStore
}

// InMemoryDataStore struct represents the in-memory data store, including data, alerts, and a file for persistence
type InMemoryDataStore struct {
	Data    map[string]Data
	Alerts  map[string][]Alerts
	File    *os.File
	FileMux sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		DataStore: &InMemoryDataStore{
			Data:   make(map[string]Data),
			Alerts: make(map[string][]Alerts),
		},
	}
}

// The Main function sets up and starts the server and also defines routes.
func main() {
	server := NewServer()
	r := chi.NewRouter()

	// Route requests
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from home"))
	})
	r.Post("/alerts", server.WriteAlert)
	r.Get("/alerts/service_id={service_id}&start_ts={alert_ts}&end_ts={alert_end_ts}", server.ReadAlerts)

	// start the Server
	log.Println("Server started...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(fmt.Sprintf("%+v", err))
	}
}

// Result struct represents the result of an operation, also includes both alert ID and error.
type Result struct {
	AlertID string
	Err     error
}

// ReqData struct represents the data structure request for writing alerts.

type ReqData struct {
	AlertID     string `json:"alert_id"`
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	Model       string `json:"model"`
	AlertType   string `json:"alert_type"`
	AlertTS     string `json:"alert_ts"`
	Severity    string `json:"severity"`
	TeamSlack   string `json:"team_slack"`
}

type Data struct {
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
}

type Alerts struct {
	AlertID   string `json:"alert_id"`
	Model     string `json:"model"`
	AlertType string `json:"alert_type"`
	AlertTs   string `json:"alert_ts"`
	Severity  string `json:"severity"`
	TeamSlack string `json:"team_slack"`
	ServiceID string `json:"service_id"`
}

// POST Request Handler (Write Alert)
func (s *Server) WriteAlert(w http.ResponseWriter, r *http.Request) {
	var data ReqData

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.DataStore.FileMux.Lock()
	defer s.DataStore.FileMux.Unlock()

	// Save the data in memory
	s.DataStore.Data[data.ServiceID] = Data{
		ServiceID:   data.ServiceID,
		ServiceName: data.ServiceName,
	}
	s.DataStore.Alerts[data.ServiceID] = append(s.DataStore.Alerts[data.ServiceID], Alerts{
		AlertID:   data.AlertID,
		Model:     data.Model,
		AlertType: data.AlertType,
		AlertTs:   data.AlertTS,
		Severity:  data.Severity,
		TeamSlack: data.TeamSlack,
		ServiceID: data.ServiceID,
	})

	// Write the data to a JSON file
	if err := s.writeDataToFile(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := Result{AlertID: data.AlertID, Err: nil}
	jsonResult, err := json.Marshal(result)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GET Request Handler (Read Alerts)
func (s *Server) ReadAlerts(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "service_id")
	alertTS := chi.URLParam(r, "alert_ts")
	alertEndTS := chi.URLParam(r, "alert_end_ts")

	s.DataStore.FileMux.RLock()
	defer s.DataStore.FileMux.RUnlock()

	alerts, ok := s.DataStore.Alerts[serviceID]
	if !ok {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	// Filter alerts by the specified time range
	var filteredAlerts []Alerts
	for _, alert := range alerts {
		if alert.AlertTs >= alertTS && alert.AlertTs <= alertEndTS {
			filteredAlerts = append(filteredAlerts, alert)
		}
	}

	if len(filteredAlerts) == 0 {
		http.Error(w, "No alerts found in the specified time range", http.StatusNotFound)
		return
	}

	result := Result{AlertID: filteredAlerts[0].AlertID, Err: nil}
	jsonResult, err := json.Marshal(result)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// writeDataToFile saves the data in memory to a JSON file
func (s *Server) writeDataToFile() error {
	file, err := os.Create("data.json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(s.DataStore); err != nil {
		return err
	}

	return nil
}
