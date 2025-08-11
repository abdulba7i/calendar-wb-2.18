package calendar

type Service struct {
	Calendar *Calendar
}

func NewService() *Service {
	return &Service{
		Calendar: NewCalendar(),
	}
}
