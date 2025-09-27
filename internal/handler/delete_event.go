package handler

import (
	"encoding/json"
	"net/http"
	"wb-calendar/internal/handler/types"
	"wb-calendar/internal/handler/validation"
	"wb-calendar/pkg/response"

	"github.com/gin-gonic/gin"
)

// DeleteEventHandler обрабатывает запросы POST /delete_event
func (h *CalendarHandler) DeleteEventHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.DeleteEventRequest

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
		if err := validation.ValidateDeleteEventRequest(&req); err != nil {
			response.JSONError(ctx, http.StatusBadRequest, err.Error())
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
