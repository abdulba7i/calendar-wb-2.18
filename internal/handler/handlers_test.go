package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"wb-calendar/internal/calendar"
	"wb-calendar/internal/handler/types"

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
			requestBody: types.CreateEventRequest{
				UserID: 1,
				Date:   "2023-12-25",
				Title:  "Christmas",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid form request",
			requestBody: types.CreateEventRequest{
				UserID: 1,
				Date:   "2023-12-25",
				Title:  "Christmas",
			},
			contentType:    "application/x-www-form-urlencoded",
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid user_id",
			requestBody: types.CreateEventRequest{
				UserID: 0,
				Date:   "2023-12-25",
				Title:  "Christmas",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty title",
			requestBody: types.CreateEventRequest{
				UserID: 1,
				Date:   "2023-12-25",
				Title:  "",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			if tt.contentType == "application/json" {
				reqBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			} else {
				reqBody = []byte("user_id=1&date=2023-12-25&title=Christmas")
			}

			req, err := http.NewRequest("POST", "/api/create_event", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestUpdateEventHandler(t *testing.T) {
	router, service := setupTestRouter()

	event, err := service.Calendar.CreateEvent(1, time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas")
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	tests := []struct {
		name           string
		requestBody    interface{}
		contentType    string
		expectedStatus int
	}{
		{
			name: "valid JSON request",
			requestBody: types.UpdateEventRequest{
				ID:    event.ID,
				Date:  "2023-12-26",
				Title: "Boxing Day",
			},
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
		},
		{
			name: "valid form request",
			requestBody: types.UpdateEventRequest{
				ID:    event.ID,
				Date:  "2023-12-26",
				Title: "Boxing Day",
			},
			contentType:    "application/x-www-form-urlencoded",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			if tt.contentType == "application/json" {
				reqBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			} else {
				reqBody = []byte(fmt.Sprintf("id=%d&date=2023-12-26&title=Boxing Day", event.ID))
			}

			req, err := http.NewRequest("POST", "/api/update_event", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestDeleteEventHandler(t *testing.T) {
	router, service := setupTestRouter()

	tests := []struct {
		name           string
		requestBody    interface{}
		contentType    string
		expectedStatus int
		setupEvent     bool
	}{
		{
			name: "valid JSON request",
			requestBody: types.DeleteEventRequest{
				ID: 1,
			},
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			setupEvent:     true,
		},
		{
			name: "valid form request",
			requestBody: types.DeleteEventRequest{
				ID: 1,
			},
			contentType:    "application/x-www-form-urlencoded",
			expectedStatus: http.StatusOK,
			setupEvent:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var eventID int
			if tt.setupEvent {
				event, err := service.Calendar.CreateEvent(1, time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas")
				if err != nil {
					t.Fatalf("Failed to create event: %v", err)
				}
				eventID = event.ID
			}

			var reqBody []byte
			var err error

			if tt.contentType == "application/json" {
				req := tt.requestBody.(types.DeleteEventRequest)
				req.ID = eventID
				reqBody, err = json.Marshal(req)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			} else {
				reqBody = []byte(fmt.Sprintf("id=%d", eventID))
			}

			req, err := http.NewRequest("POST", "/api/delete_event", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetEventsForDayHandler(t *testing.T) {
	router, service := setupTestRouter()

	_, err := service.Calendar.CreateEvent(1, time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas")
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	_, err = service.Calendar.CreateEvent(1, time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas Eve")
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name: "valid request",
			requestBody: types.GetEventsRequest{
				UserID: 1,
				Date:   "2023-12-25",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid user_id",
			requestBody: types.GetEventsRequest{
				UserID: 0,
				Date:   "2023-12-25",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty date",
			requestBody: types.GetEventsRequest{
				UserID: 1,
				Date:   "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest("GET", "/api/events_for_day", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetEventsForWeekHandler(t *testing.T) {
	router, service := setupTestRouter()

	_, err := service.Calendar.CreateEvent(1, time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas")
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name: "valid request",
			requestBody: types.GetEventsRequest{
				UserID: 1,
				Date:   "2023-12-25",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest("GET", "/api/events_for_week", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetEventsForMonthHandler(t *testing.T) {
	router, service := setupTestRouter()

	_, err := service.Calendar.CreateEvent(1, time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC), "Christmas")
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name: "valid request",
			requestBody: types.GetEventsRequest{
				UserID: 1,
				Date:   "2023-12-25",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest("GET", "/api/events_for_month", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
