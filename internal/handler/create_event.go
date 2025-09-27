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

// CreateEventHandler обрабатывает запросы POST /create_event
func (h *CalendarHandler) CreateEventHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.CreateEventRequest

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
		if err := validation.ValidateCreateEventRequest(&req); err != nil {
			response.JSONError(ctx, http.StatusBadRequest, err.Error())
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
