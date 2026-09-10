package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/instituto-libcom/go-api-core/pagination"
)

func TestSuccessResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	data := map[string]string{"id": "123", "name": "Test"}

	Success(recorder, http.StatusOK, data)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", recorder.Code)
	}
	
	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type application/json, got %s", contentType)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, `"success":true`) {
		t.Errorf("Expected success to be true, got %s", body)
	}
	if !strings.Contains(body, `"data":{"id":"123","name":"Test"}`) {
		t.Errorf("Expected data to be serialized correctly, got %s", body)
	}
	if !strings.Contains(body, `"error":null`) {
		t.Errorf("Expected error to be null, got %s", body)
	}
}

func TestPaginatedResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	items := []map[string]string{{"id": "1"}, {"id": "2"}}
	meta := pagination.Meta{
		Page:          0,
		Size:          25,
		TotalElements: 10,
		TotalPages:    1,
	}

	Paginated(recorder, http.StatusOK, items, meta)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", recorder.Code)
	}

	var parsed struct {
		Success bool `json:"success"`
		Data    struct {
			Items []map[string]string `json:"items"`
			Meta  pagination.Meta     `json:"meta"`
		} `json:"data"`
		Error *ErrorDetails `json:"error"`
	}

	err := json.Unmarshal(recorder.Body.Bytes(), &parsed)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if !parsed.Success {
		t.Errorf("Expected success to be true")
	}
	if parsed.Error != nil {
		t.Errorf("Expected error to be nil")
	}
	if len(parsed.Data.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(parsed.Data.Items))
	}
	if parsed.Data.Meta.TotalElements != 10 {
		t.Errorf("Expected total_elements 10, got %d", parsed.Data.Meta.TotalElements)
	}
	if parsed.Data.Meta.TotalPages != 1 {
		t.Errorf("Expected total_pages 1, got %d", parsed.Data.Meta.TotalPages)
	}
	
	// ensure specific string match for JSON contract
	bodyStr := recorder.Body.String()
	if !strings.Contains(bodyStr, `"total_elements":10`) || !strings.Contains(bodyStr, `"total_pages":1`) {
		t.Errorf("JSON keys for pagination meta do not match expected contract. Body: %s", bodyStr)
	}
}

func TestErrorResponse(t *testing.T) {
	recorder := httptest.NewRecorder()

	Error(recorder, http.StatusBadRequest, "INVALID_INPUT", "Input validation failed")

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code 400, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	if !strings.Contains(body, `"success":false`) {
		t.Errorf("Expected success to be false, got %s", body)
	}
	if !strings.Contains(body, `"data":null`) {
		t.Errorf("Expected data to be null, got %s", body)
	}
	if !strings.Contains(body, `"code":"INVALID_INPUT"`) {
		t.Errorf("Expected code to be serialized correctly, got %s", body)
	}
	if !strings.Contains(body, `"message":"Input validation failed"`) {
		t.Errorf("Expected message to be serialized correctly, got %s", body)
	}
}
