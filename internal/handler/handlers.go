package handler

import (
	"wb-calendar/internal/calendar"
)

// CalendarHandler обрабатывает HTTP-запросы для операций с календарем
type CalendarHandler struct {
	service calendar.Service
}

// NewCalendarHandler создает новый экземпляр CalendarHandler
func NewCalendarHandler(service calendar.Service) *CalendarHandler {
	return &CalendarHandler{service: service}
}
