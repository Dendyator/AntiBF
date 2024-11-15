package api_test

import (
	"bytes"
	"encoding/json"
	"github.com/Dendyator/AntiBF/internal/db"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Dendyator/AntiBF/api"
	"github.com/Dendyator/AntiBF/internal/config" //nolint
	"github.com/Dendyator/AntiBF/internal/core"   //nolint
	"github.com/Dendyator/AntiBF/internal/logger" //nolint
)

func TestMain(m *testing.M) {
	log := logger.New("debug")
	db.InitRedis("localhost:6379", log) // убедитесь, что адрес Redis правильный
	defer db.CloseRedis()
	code := m.Run()
	os.Exit(code)
}

func TestHandleAuth(t *testing.T) {
	mockLogger := logger.New("info")
	core.InitLogger(mockLogger)
	api.InitLogger(mockLogger)

	core.InitRateLimiter(config.RateLimiterConfig{
		LoginLimit:    10,
		PasswordLimit: 100,
		IPLimit:       1000,
	})

	core.WhitelistFunc = func(ip string) bool {
		return false
	}
	core.BlacklistFunc = func(ip string) bool {
		return false
	}

	reqBody, _ := json.Marshal(map[string]string{
		"login":    "testuser",
		"password": "testpass",
		"ip":       "192.168.1.30/25",
	})

	req, err := http.NewRequest("POST", "/auth", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.HandleAuth)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := map[string]bool{"ok": true}
	var response map[string]bool
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse JSON response: %v", err)
	}

	if response["ok"] != expected["ok"] {
		t.Errorf("Handler returned unexpected body: got %v want %v", response, expected)
	}
}

func TestHandleManageList(t *testing.T) {
	mockLogger := logger.New("info")
	core.InitLogger(mockLogger)
	api.InitLogger(mockLogger)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"listType": "white",
		"subnet":   "192.168.1.0/24",
		"add":      true,
	})

	req, err := http.NewRequest("POST", "/manage_list", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.HandleManageList)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]bool
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	expected := true
	if response["success"] != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", response["success"], expected)
	}
}

func TestHandleCheckList(t *testing.T) {
	mockLogger := logger.New("info")
	core.InitLogger(mockLogger)
	api.InitLogger(mockLogger)

	requestBody, _ := json.Marshal(map[string]string{
		"listType": "black",
		"subnet":   "192.168.1.1/25",
	})

	req, err := http.NewRequest("POST", "/check_list", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.HandleCheckList)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]bool
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	expected := false
	if response["in_list"] != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", response["in_list"], expected)
	}
}
