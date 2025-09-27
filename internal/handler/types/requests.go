package types

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

// GetEventsRequest структура для получения событий
type GetEventsRequest struct {
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
}
