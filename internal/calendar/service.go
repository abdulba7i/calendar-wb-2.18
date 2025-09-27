package calendar

// Service предоставляет бизнес-логику календаря
type Service struct {
	Calendar *Calendar
}

// NewService создает новый экземпляр
func NewService() *Service {
	return &Service{
		Calendar: NewCalendar(),
	}
}
