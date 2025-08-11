package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"wb-calendar/internal/calendar"
	"wb-calendar/pkg/response"

	"github.com/gin-gonic/gin"
)

type CalendarHandler struct {
	service calendar.Service
}

func NewCalendarHandler(service calendar.Service) *CalendarHandler {
	return &CalendarHandler{service: service}
}

// CreateEventRequest структура для создания события
type CreateEventRequest struct {
	UserID int    `json:"user_id" form:"user_id"`
	Date   string `json:"date" form:"date"`
	Title  string `json:"title" form:"title"`
}

// UpdateEventRequest структура для обновления события
type UpdateEventRequest struct {
	ID    int    `json:"id" form:"id"`
	Date  string `json:"date" form:"date"`
	Title string `json:"title" form:"title"`
}

// DeleteEventRequest структура для удаления события
type DeleteEventRequest struct {
	ID int `json:"id" form:"id"`
}

func (h *CalendarHandler) CreateEventHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CreateEventRequest

		// Поддерживаем оба формата: JSON и form
		contentType := ctx.GetHeader("Content-Type")
		if contentType == "application/json" {
			if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
				response.JSONError(ctx, http.StatusBadRequest, "invalid JSON request body")
				return
			}
		} else {
			if err := ctx.ShouldBind(&req); err != nil {
				response.JSONError(ctx, http.StatusBadRequest, "invalid form data")
				return
			}
		}

		// Валидация
		if req.UserID <= 0 {
			response.JSONError(ctx, http.StatusBadRequest, "user_id must be positive")
			return
		}
		if req.Title == "" {
			response.JSONError(ctx, http.StatusBadRequest, "title cannot be empty")
			return
		}

		// Парсинг даты
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			response.JSONError(ctx, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
			return
		}

		event, err := h.service.Calendar.CreateEvent(req.UserID, date, req.Title)
		if err != nil {
			response.JSONError(ctx, http.StatusInternalServerError, "failed to create event")
			return
		}

		response.JSONResult(ctx, event)
	}
}

func (h *CalendarHandler) UpdateEventHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req UpdateEventRequest

		// Поддерживаем оба формата: JSON и form
		contentType := ctx.GetHeader("Content-Type")
		if contentType == "application/json" {
			if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
				response.JSONError(ctx, http.StatusBadRequest, "invalid JSON request body")
				return
			}
		} else {
			if err := ctx.ShouldBind(&req); err != nil {
				response.JSONError(ctx, http.StatusBadRequest, "invalid form data")
				return
			}
		}

		// Валидация
		if req.ID <= 0 {
			response.JSONError(ctx, http.StatusBadRequest, "id must be positive")
			return
		}
		if req.Title == "" {
			response.JSONError(ctx, http.StatusBadRequest, "title cannot be empty")
			return
		}

		// Парсинг даты
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			response.JSONError(ctx, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
			return
		}

		if err := h.service.Calendar.UpdateEvent(req.ID, date, req.Title); err != nil {
			if err.Error() == "event not found" {
				response.JSONError(ctx, http.StatusServiceUnavailable, "event not found")
				return
			}
			response.JSONError(ctx, http.StatusInternalServerError, "failed to update event")
			return
		}

		response.JSONResult(ctx, "event updated successfully")
	}
}

func (h *CalendarHandler) DeleteEventHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req DeleteEventRequest

		contentType := ctx.GetHeader("Content-Type")
		if contentType == "application/json" {
			if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
				response.JSONError(ctx, http.StatusBadRequest, "invalid JSON request body")
				return
			}
		} else {
			if err := ctx.ShouldBind(&req); err != nil {
				response.JSONError(ctx, http.StatusBadRequest, "invalid form data")
				return
			}
		}

		// Валидация
		if req.ID <= 0 {
			response.JSONError(ctx, http.StatusBadRequest, "id must be positive")
			return
		}

		if err := h.service.Calendar.DeleteEvent(req.ID); err != nil {
			if err.Error() == "event not found" {
				response.JSONError(ctx, http.StatusServiceUnavailable, "event not found")
				return
			}
			response.JSONError(ctx, http.StatusInternalServerError, "failed to delete event")
			return
		}

		response.JSONResult(ctx, "event deleted successfully")
	}
}

func (h *CalendarHandler) GetEventsForDayHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req struct {
			UserID int    `json:"user_id"`
			Date   string `json:"date"`
		}

		if err := ctx.BindJSON(&req); err != nil {
			response.JSONError(ctx, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.UserID <= 0 {
			response.JSONError(ctx, http.StatusBadRequest, "invalid user_id")
			return
		}

		if req.Date == "" {
			response.JSONError(ctx, http.StatusBadRequest, "date parameter is required")
			return
		}

		day, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			response.JSONError(ctx, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
			return
		}

		events := h.service.Calendar.GetEventsForDay(req.UserID, day)
		response.JSONResult(ctx, events)
	}
}

func (h *CalendarHandler) GetEventsForWeekHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req struct {
			UserID int    `json:"user_id"`
			Date   string `json:"date"`
		}

		if err := ctx.BindJSON(&req); err != nil {
			response.JSONError(ctx, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.UserID <= 0 {
			response.JSONError(ctx, http.StatusBadRequest, "invalid user_id")
			return
		}

		if req.Date == "" {
			response.JSONError(ctx, http.StatusBadRequest, "date parameter is required")
			return
		}

		day, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			response.JSONError(ctx, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
			return
		}

		events := h.service.Calendar.GetEventsForWeek(req.UserID, day)
		response.JSONResult(ctx, events)
	}
}

func (h *CalendarHandler) GetEventsForMonthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req struct {
			UserID int    `json:"user_id"`
			Date   string `json:"date"`
		}

		if err := ctx.BindJSON(&req); err != nil {
			response.JSONError(ctx, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.UserID <= 0 {
			response.JSONError(ctx, http.StatusBadRequest, "invalid user_id")
			return
		}

		if req.Date == "" {
			response.JSONError(ctx, http.StatusBadRequest, "date parameter is required")
			return
		}

		day, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			response.JSONError(ctx, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
			return
		}

		events := h.service.Calendar.GetEventsForMonth(req.UserID, day)
		response.JSONResult(ctx, events)
	}
}
