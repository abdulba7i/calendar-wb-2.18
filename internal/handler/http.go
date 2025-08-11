package handler

import (
	"wb-calendar/internal/calendar"
	"wb-calendar/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoute(service *calendar.Service) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.LoggingMiddleware())

	calendarHandler := NewCalendarHandler(*service)

	r.POST("/create_event", calendarHandler.CreateEventHandler())
	r.POST("/update_event", calendarHandler.UpdateEventHandler())
	r.POST("/delete_event", calendarHandler.DeleteEventHandler())

	r.GET("/events_for_day", calendarHandler.GetEventsForDayHandler())
	r.GET("/events_for_week", calendarHandler.GetEventsForWeekHandler())
	r.GET("/events_for_month", calendarHandler.GetEventsForMonthHandler())

	return r
}
