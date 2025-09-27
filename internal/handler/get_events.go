package handler

import (
	"net/http"
	"time"
	"wb-calendar/internal/handler/types"
	"wb-calendar/internal/handler/validation"
	"wb-calendar/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetEventsForDayHandler обрабатывает запросы GET /events_for_day
func (h *CalendarHandler) GetEventsForDayHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.GetEventsRequest

		if err := ctx.BindJSON(&req); err != nil {
			response.JSONError(ctx, http.StatusBadRequest, "invalid request body")
			return
		}

		// Валидация
		if err := validation.ValidateGetEventsRequest(&req); err != nil {
			response.JSONError(ctx, http.StatusBadRequest, err.Error())
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

// GetEventsForWeekHandler handles GET /events_for_week requests
func (h *CalendarHandler) GetEventsForWeekHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.GetEventsRequest

		if err := ctx.BindJSON(&req); err != nil {
			response.JSONError(ctx, http.StatusBadRequest, "invalid request body")
			return
		}

		// Валидация
		if err := validation.ValidateGetEventsRequest(&req); err != nil {
			response.JSONError(ctx, http.StatusBadRequest, err.Error())
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

// GetEventsForMonthHandler handles GET /events_for_month requests
func (h *CalendarHandler) GetEventsForMonthHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.GetEventsRequest

		if err := ctx.BindJSON(&req); err != nil {
			response.JSONError(ctx, http.StatusBadRequest, "invalid request body")
			return
		}

		// Валидация
		if err := validation.ValidateGetEventsRequest(&req); err != nil {
			response.JSONError(ctx, http.StatusBadRequest, err.Error())
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
