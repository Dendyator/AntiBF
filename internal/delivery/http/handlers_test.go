package http_test

import (
	"bytes"
	"encoding/json"
	api2 "github.com/Dendyator/AntiBF/internal/delivery/http"
	"github.com/Dendyator/AntiBF/internal/repositories"
	"github.com/Dendyator/AntiBF/pkg/config"
	"github.com/Dendyator/AntiBF/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dendyator/AntiBF/internal/usecase" //nolint
)

func TestHandleAuth(t *testing.T) {
	mockLogger := logger.New("info")

	// Создаем моковые репозитории
	mockRepo := &repositories.MockRepository{}
	mockUserRepo := &repositories.MockUserRepository{}

	// Создаем RateLimiter
	rateLimiter := usecase.NewRateLimiter(mockRepo, mockUserRepo, config.RateLimiterConfig{
		LoginLimit:    10,
		PasswordLimit: 100,
		IPLimit:       1000,
	}, mockLogger)

	// Подготавливаем запрос
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

	// Используем замыкание для передачи rateLimiter в обработчик
	handler := http.HandlerFunc(api2.HandleAuth(rateLimiter))

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем ответ
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

	// Создаем моковые репозитории
	mockRepo := &repositories.MockRepository{}
	mockUserRepo := &repositories.MockUserRepository{}

	// Создаем RateLimiter
	rateLimiter := usecase.NewRateLimiter(mockRepo, mockUserRepo, config.RateLimiterConfig{
		LoginLimit:    10,
		PasswordLimit: 100,
		IPLimit:       1000,
	}, mockLogger)

	// Подготавливаем запрос для добавления в whitelist
	requestBody, _ := json.Marshal(map[string]interface{}{
		"listType": "white",
		"subnet":   "192.168.1.0/24",
		"add":      true,
	})
	req, err := http.NewRequest("POST", "/manage_list", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Создаем обработчик с переданным rateLimiter
	handler := api2.HandleManageList(rateLimiter)

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем ответ
	var response map[string]bool
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse JSON response: %v", err)
	}
	expected := map[string]bool{"success": true}
	if response["success"] != expected["success"] {
		t.Errorf("Handler returned unexpected body: got %v want %v", response, expected)
	}

	// Тест с некорректным форматом subnet
	invalidRequestBody, _ := json.Marshal(map[string]interface{}{
		"listType": "white",
		"subnet":   "invalid-subnet",
		"add":      true,
	})
	invalidReq, err := http.NewRequest("POST", "/manage_list", bytes.NewBuffer(invalidRequestBody))
	if err != nil {
		t.Fatal(err)
	}

	invalidRR := httptest.NewRecorder()
	handler.ServeHTTP(invalidRR, invalidReq)

	// Проверяем статус код для ошибки
	if status := invalidRR.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for invalid subnet: got %v want %v", status, http.StatusBadRequest)
	}

	// Тест с некорректным типом списка
	invalidListTypeRequestBody, _ := json.Marshal(map[string]interface{}{
		"listType": "invalid",
		"subnet":   "192.168.1.0/24",
		"add":      true,
	})
	invalidListTypeReq, err := http.NewRequest("POST", "/manage_list", bytes.NewBuffer(invalidListTypeRequestBody))
	if err != nil {
		t.Fatal(err)
	}

	invalidListTypeRR := httptest.NewRecorder()
	handler.ServeHTTP(invalidListTypeRR, invalidListTypeReq)

	// Проверяем статус код для ошибки типа списка
	if status := invalidListTypeRR.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for invalid list type: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestHandleCheckList(t *testing.T) {
	mockLogger := logger.New("info")

	// Создаем моковые репозитории
	mockRepo := &repositories.MockRepository{}
	mockUserRepo := &repositories.MockUserRepository{}

	// Создаем RateLimiter
	rateLimiter := usecase.NewRateLimiter(mockRepo, mockUserRepo, config.RateLimiterConfig{
		LoginLimit:    10,
		PasswordLimit: 100,
		IPLimit:       1000,
	}, mockLogger)

	// Подготавливаем запрос для проверки в whitelist
	requestBody, _ := json.Marshal(map[string]string{
		"listType": "white",
		"subnet":   "192.168.1.0/24",
	})
	req, err := http.NewRequest("POST", "/check_list", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Создаем обработчик с переданным rateLimiter
	handler := api2.HandleCheckList(rateLimiter)

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем ответ
	var response map[string]bool
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Could not parse JSON response: %v", err)
	}
	expected := map[string]bool{"in_list": false} // По умолчанию subnet не должен быть в списке
	if response["in_list"] != expected["in_list"] {
		t.Errorf("Handler returned unexpected body: got %v want %v", response, expected)
	}

	// Тест с некорректным форматом subnet
	invalidRequestBody, _ := json.Marshal(map[string]string{
		"listType": "white",
		"subnet":   "invalid-subnet",
	})
	invalidReq, err := http.NewRequest("POST", "/check_list", bytes.NewBuffer(invalidRequestBody))
	if err != nil {
		t.Fatal(err)
	}

	invalidRR := httptest.NewRecorder()
	handler.ServeHTTP(invalidRR, invalidReq)

	// Проверяем статус код для ошибки
	if status := invalidRR.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for invalid subnet: got %v want %v", status, http.StatusBadRequest)
	}

	// Тест с некорректным типом списка
	invalidListTypeRequestBody, _ := json.Marshal(map[string]string{
		"listType": "invalid",
		"subnet":   "192.168.1.0/24",
	})
	invalidListTypeReq, err := http.NewRequest("POST", "/check_list", bytes.NewBuffer(invalidListTypeRequestBody))
	if err != nil {
		t.Fatal(err)
	}

	invalidListTypeRR := httptest.NewRecorder()
	handler.ServeHTTP(invalidListTypeRR, invalidListTypeReq)

	// Проверяем статус код для ошибки типа списка
	if status := invalidListTypeRR.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code for invalid list type: got %v want %v", status, http.StatusBadRequest)
	}
}
