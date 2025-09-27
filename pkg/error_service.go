package pkg

import "errors"

// ErrEventNotFound возвращается, если событие не найдено
var (
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidDate   = errors.New("invalid date format")
)
