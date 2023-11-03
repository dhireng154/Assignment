package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestWriteAlert(t *testing.T) {
	// Create a new server instance
	server := NewServer()

	// Create a test request payload
	alertData := ReqData{
		AlertID:     "b950482e9911ec7e41f7ca5e5d9a424f",
		ServiceID:   "my_test_service_id",
		ServiceName: "my_test_service",
		Model:       "TestModel",
		AlertType:   "Critical",
		AlertTS:     time.Now().UTC().Format(time.RFC3339),
		Severity:    "High",
		TeamSlack:   "testteam",
	}

	// Convert the payload to JSON
	jsonData, err := json.Marshal(alertData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON data: %v", err)
	}

	// Create a test request with the JSON payload
	req, err := http.NewRequest("POST", "/alerts", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Create a test HTTP response recorder
	recorder := httptest.NewRecorder()

	// Call the WriteAlert function with the test request and recorder
	server.WriteAlert(recorder, req)

	// Check the status code of the response
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// Decode the response body into a Result struct
	var result Result
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	// Check if the AlertID in the response matches the expected AlertID
	expectedAlertID := "b950482e9911ec7e41f7ca5e5d9a424f"
	if result.AlertID != expectedAlertID {
		t.Errorf("Expected AlertID %s, got %s", expectedAlertID, result.AlertID)
	}

	// Check if the Error in the response is nil
	if result.Err != nil {
		t.Errorf("Expected error to be nil, got %v", result.Err)
	}
}

func TestReadAlerts(t *testing.T) {
	// Create a new server instance
	server := NewServer()

	// Populate the in-memory data store with test data
	testServiceID := "my_test_service_id"
	testAlert := Alerts{
		AlertID:   "b950482e9911ec7e41f7ca5e5d9a424f",
		Model:     "TestModel",
		AlertType: "Critical",
		AlertTs:   time.Now().UTC().Format(time.RFC3339),
		Severity:  "High",
		TeamSlack: "testteam",
		ServiceID: testServiceID,
	}
	server.DataStore.Alerts[testServiceID] = []Alerts{testAlert}

	// Create a test request with the required URL parameters
	url := fmt.Sprintf("/alerts/service_id=%s&start_ts=%s&end_ts=%s",
		testServiceID, time.Now().Add(-time.Hour).UTC().Format(time.RFC3339),
		time.Now().Add(time.Hour).UTC().Format(time.RFC3339))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Create a test HTTP response recorder
	recorder := httptest.NewRecorder()

	// Call the ReadAlerts function with the test request and recorder
	server.ReadAlerts(recorder, req)

	// Check the status code of the response
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, recorder.Code)
	}

	// Decode the response body into a Result struct
	var result Result
	err = json.Unmarshal(recorder.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	// Check if the AlertID in the response matches the expected AlertID
	expectedAlertID := testAlert.AlertID
	if result.AlertID != expectedAlertID {
		t.Errorf("Expected AlertID %s, got %s", expectedAlertID, result.AlertID)
	}

	// Check if the Error in the response is nil
	if result.Err != nil {
		t.Errorf("Expected error to be nil, got %v", result.Err)
	}
}

func TestWriteDataToFile(t *testing.T) {
	// Create a new server instance
	server := NewServer()

	// Create a test file for writing
	testFileName := "test_data.json"
	server.DataStore.File, _ = os.Create(testFileName)

	// Call the writeDataToFile function
	err := server.writeDataToFile()
	if err != nil {
		t.Fatalf("writeDataToFile returned an error: %v", err)
	}

	// Check if the test file exists
	_, err = os.Stat(testFileName)
	if os.IsNotExist(err) {
		t.Errorf("writeDataToFile did not create the expected file: %v", err)
	}

	// Closing the file.
	server.DataStore.File.Close()

	// Clean up: Remove the test file
	err = os.Remove(testFileName)
	if err != nil {
		t.Fatalf("Failed to remove test file: %v", err)
	}
}
