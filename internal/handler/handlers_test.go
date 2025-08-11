package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"wb-calendar/internal/calendar"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() (*gin.Engine, *calendar.Service) {
	gin.SetMode(gin.TestMode)
	service := calendar.NewService()
	handler := NewCalendarHandler(*service)
	router := gin.New()

	api := router.Group("/api")
	{
		api.POST("/create_event", handler.CreateEventHandler())
		api.POST("/update_event", handler.UpdateEventHandler())
		api.POST("/delete_event", handler.DeleteEventHandler())
		api.GET("/events_for_day", handler.GetEventsForDayHandler())
		api.GET("/events_for_week", handler.GetEventsForWeekHandler())
		api.GET("/events_for_month", handler.GetEventsForMonthHandler())
	}

	return router, service
}

func TestCreateEventHandler(t *testing.T) {
	router, _ := setupTestRouter()

	tests := []struct {
		name           string
		requestBody    interface{}
		contentType    string
		expectedStatus int
	}{
		{
			name: "valid JSON request",
			requestBody: CreateEventRequest{
				UserID: 1,
				Date:   "2023-12-25",
				Title:  "Christmas",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid date format",
			requestBody: CreateEventRequest{
				UserID: 1,
				Date:   "invalid-date",
				Title:  "Christmas",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty title",
			requestBody: CreateEventRequest{
				UserID: 1,
				Date:   "2023-12-25",
				Title:  "",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid user_id",
			requestBody: CreateEventRequest{
				UserID: 0,
				Date:   "2023-12-25",
				Title:  "Christmas",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/create_event", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestUpdateEventHandler(t *testing.T) {
	router, service := setupTestRouter()

	event, _ := service.Calendar.CreateEvent(1, time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas")

	tests := []struct {
		name           string
		requestBody    interface{}
		contentType    string
		expectedStatus int
	}{
		{
			name: "valid update request",
			requestBody: UpdateEventRequest{
				ID:    event.ID,
				Date:  "2023-12-26",
				Title: "Boxing Day",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
		},
		{
			name: "event not found",
			requestBody: UpdateEventRequest{
				ID:    999,
				Date:  "2023-12-26",
				Title: "Boxing Day",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/update_event", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestDeleteEventHandler(t *testing.T) {
	router, service := setupTestRouter()

	event, _ := service.Calendar.CreateEvent(1, time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas")

	tests := []struct {
		name           string
		requestBody    interface{}
		contentType    string
		expectedStatus int
	}{
		{
			name: "valid delete request",
			requestBody: DeleteEventRequest{
				ID: event.ID,
			},
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
		},
		{
			name: "event not found",
			requestBody: DeleteEventRequest{
				ID: 999,
			},
			contentType:    "application/json",
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/delete_event", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetEventsForDayHandler(t *testing.T) {
	router, service := setupTestRouter()

	date1 := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	service.Calendar.CreateEvent(1, date1, "Christmas")

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid request",
			requestBody: map[string]interface{}{
				"user_id": 1,
				"date":    "2023-12-25",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing user_id",
			requestBody:    map[string]interface{}{"date": "2023-12-25"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing date",
			requestBody:    map[string]interface{}{"user_id": 1},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid user_id",
			requestBody: map[string]interface{}{
				"user_id": "abc",
				"date":    "2023-12-25",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid date format",
			requestBody: map[string]interface{}{
				"user_id": 1,
				"date":    "invalid-date",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("GET", "/api/events_for_day", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetEventsForWeekHandler(t *testing.T) {
	router, service := setupTestRouter()

	date1 := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	service.Calendar.CreateEvent(1, date1, "Christmas")

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid request",
			requestBody: map[string]interface{}{
				"user_id": 1,
				"date":    "2023-12-25",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid user_id",
			requestBody: map[string]interface{}{
				"user_id": 0,
				"date":    "2023-12-25",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("GET", "/api/events_for_week", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetEventsForMonthHandler(t *testing.T) {
	router, service := setupTestRouter()

	date1 := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	service.Calendar.CreateEvent(1, date1, "Christmas")

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "valid request",
			requestBody: map[string]interface{}{
				"user_id": 1,
				"date":    "2023-12-25",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid user_id",
			requestBody: map[string]interface{}{
				"user_id": -1,
				"date":    "2023-12-25",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("GET", "/api/events_for_month", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
