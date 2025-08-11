package calendar

import (
	"sync"
	"time"
	"wb-calendar/pkg"
)

type Calendar struct {
	events map[int]Event
	nextID int
	mutex  sync.RWMutex
}

func NewCalendar() *Calendar {
	return &Calendar{
		events: make(map[int]Event),
		nextID: 1,
		mutex:  sync.RWMutex{},
	}
}

// CreateEvent создает новое событие
func (c *Calendar) CreateEvent(userID int, date time.Time, title string) (Event, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	event := Event{
		ID:     c.nextID,
		UserID: userID,
		Date:   date,
		Title:  title,
	}

	c.events[event.ID] = event
	c.nextID++

	return event, nil
}

// UpdateEvent обновляет существующее событие
func (c *Calendar) UpdateEvent(id int, date time.Time, title string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	event, exists := c.events[id]
	if !exists {
		return pkg.ErrEventNotFound
	}

	event.Date = date
	event.Title = title

	c.events[id] = event

	return nil
}

// DeleteEvent удаляет событие
func (c *Calendar) DeleteEvent(id int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _, exists := c.events[id]; !exists {
		return pkg.ErrEventNotFound
	}

	delete(c.events, id)
	return nil
}

// GetEventsForDay возвращает события на день
func (c *Calendar) GetEventsForDay(userID int, day time.Time) []Event {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []Event

	for _, event := range c.events {
		if event.UserID == userID && isSameDay(event.Date, day) {
			result = append(result, event)
		}
	}

	return result
}

// GetEventsForWeek возвращает события на неделю
func (c *Calendar) GetEventsForWeek(userID int, day time.Time) []Event {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []Event
	year, week := day.ISOWeek()

	for _, event := range c.events {
		y, w := event.Date.ISOWeek()
		if event.UserID == userID && y == year && w == week {
			result = append(result, event)
		}
	}

	return result
}

// GetEventsForMonth возвращает события на месяц
func (c *Calendar) GetEventsForMonth(userID int, day time.Time) []Event {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var result []Event
	for _, event := range c.events {
		if event.UserID == userID && event.Date.Year() == day.Year() && event.Date.Month() == day.Month() {
			result = append(result, event)
		}
	}

	return result
}

// isSameDay проверяет, что две даты относятся к одному дню
func isSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()
}
