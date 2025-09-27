package validation

import (
	"time"
	"wb-calendar/internal/handler/types"
)

// ValidateCreateEventRequest валидирует запрос на создание события
func ValidateCreateEventRequest(req *types.CreateEventRequest) error {
	if req.UserID <= 0 {
		return &Error{Field: "user_id", Message: "user_id must be positive"}
	}
	if req.Title == "" {
		return &Error{Field: "title", Message: "title cannot be empty"}
	}
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		return &Error{Field: "date", Message: "invalid date format, expected YYYY-MM-DD"}
	}
	return nil
}

// ValidateUpdateEventRequest валидирует запрос на обновление события
func ValidateUpdateEventRequest(req *types.UpdateEventRequest) error {
	if req.ID <= 0 {
		return &Error{Field: "id", Message: "id must be positive"}
	}
	if req.Title == "" {
		return &Error{Field: "title", Message: "title cannot be empty"}
	}
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		return &Error{Field: "date", Message: "invalid date format, expected YYYY-MM-DD"}
	}
	return nil
}

// ValidateDeleteEventRequest валидирует запрос на удаление события
func ValidateDeleteEventRequest(req *types.DeleteEventRequest) error {
	if req.ID <= 0 {
		return &Error{Field: "id", Message: "id must be positive"}
	}
	return nil
}

// ValidateGetEventsRequest валидирует запрос на получение событий
func ValidateGetEventsRequest(req *types.GetEventsRequest) error {
	if req.UserID <= 0 {
		return &Error{Field: "user_id", Message: "invalid user_id"}
	}
	if req.Date == "" {
		return &Error{Field: "date", Message: "date parameter is required"}
	}
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		return &Error{Field: "date", Message: "invalid date format, expected YYYY-MM-DD"}
	}
	return nil
}

// Error представляет собой ошибку проверки
type Error struct {
	Field   string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}
