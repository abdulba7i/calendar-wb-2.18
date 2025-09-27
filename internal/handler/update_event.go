package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"wb-calendar/internal/handler/types"
	"wb-calendar/internal/handler/validation"
	"wb-calendar/pkg/response"

	"github.com/gin-gonic/gin"
)

// UpdateEventHandler обрабатывает запросы POST /update_event
func (h *CalendarHandler) UpdateEventHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UpdateEventRequest

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
		if err := validation.ValidateUpdateEventRequest(&req); err != nil {
			response.JSONError(ctx, http.StatusBadRequest, err.Error())
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
