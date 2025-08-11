package calendar

import (
	"testing"
	"time"
)

func TestNewCalendar(t *testing.T) {
	cal := NewCalendar()
	if cal == nil {
		t.Fatal("NewCalendar() returned nil")
	}
	if cal.events == nil {
		t.Fatal("events map is nil")
	}
	if cal.nextID != 1 {
		t.Fatalf("expected nextID to be 1, got %d", cal.nextID)
	}
}

func TestCreateEvent(t *testing.T) {
	cal := NewCalendar()
	userID := 1
	date := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	title := "Christmas"

	event, err := cal.CreateEvent(userID, date, title)
	if err != nil {
		t.Fatalf("CreateEvent failed: %v", err)
	}

	if event.ID != 1 {
		t.Fatalf("expected event ID to be 1, got %d", event.ID)
	}
	if event.UserID != userID {
		t.Fatalf("expected user ID to be %d, got %d", userID, event.UserID)
	}
	if !event.Date.Equal(date) {
		t.Fatalf("expected date to be %v, got %v", date, event.Date)
	}
	if event.Title != title {
		t.Fatalf("expected title to be %s, got %s", title, event.Title)
	}

	// Проверяем, что событие сохранено в календаре
	if len(cal.events) != 1 {
		t.Fatalf("expected 1 event in calendar, got %d", len(cal.events))
	}

	// Проверяем, что nextID увеличился
	if cal.nextID != 2 {
		t.Fatalf("expected nextID to be 2, got %d", cal.nextID)
	}
}

func TestUpdateEvent(t *testing.T) {
	cal := NewCalendar()
	userID := 1
	date := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	title := "Christmas"

	// Создаем событие
	event, _ := cal.CreateEvent(userID, date, title)

	// Обновляем событие
	newDate := time.Date(2023, 12, 26, 0, 0, 0, 0, time.UTC)
	newTitle := "Boxing Day"
	err := cal.UpdateEvent(event.ID, newDate, newTitle)
	if err != nil {
		t.Fatalf("UpdateEvent failed: %v", err)
	}

	// Проверяем, что событие обновлено
	updatedEvent := cal.events[event.ID]
	if !updatedEvent.Date.Equal(newDate) {
		t.Fatalf("expected updated date to be %v, got %v", newDate, updatedEvent.Date)
	}
	if updatedEvent.Title != newTitle {
		t.Fatalf("expected updated title to be %s, got %s", newTitle, updatedEvent.Title)
	}
}

func TestUpdateEventNotFound(t *testing.T) {
	cal := NewCalendar()
	date := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	title := "Christmas"

	// Пытаемся обновить несуществующее событие
	err := cal.UpdateEvent(999, date, title)
	if err == nil {
		t.Fatal("expected error when updating non-existent event")
	}
}

func TestDeleteEvent(t *testing.T) {
	cal := NewCalendar()
	userID := 1
	date := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	title := "Christmas"

	// Создаем событие
	event, _ := cal.CreateEvent(userID, date, title)

	// Проверяем, что событие существует
	if len(cal.events) != 1 {
		t.Fatalf("expected 1 event before deletion, got %d", len(cal.events))
	}

	// Удаляем событие
	err := cal.DeleteEvent(event.ID)
	if err != nil {
		t.Fatalf("DeleteEvent failed: %v", err)
	}

	// Проверяем, что событие удалено
	if len(cal.events) != 0 {
		t.Fatalf("expected 0 events after deletion, got %d", len(cal.events))
	}
}

func TestDeleteEventNotFound(t *testing.T) {
	cal := NewCalendar()

	// Пытаемся удалить несуществующее событие
	err := cal.DeleteEvent(999)
	if err == nil {
		t.Fatal("expected error when deleting non-existent event")
	}
}

func TestGetEventsForDay(t *testing.T) {
	cal := NewCalendar()
	userID := 1
	date1 := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2023, 12, 26, 0, 0, 0, 0, time.UTC)

	// Создаем два события на разные дни
	cal.CreateEvent(userID, date1, "Christmas")
	cal.CreateEvent(userID, date2, "Boxing Day")

	// Получаем события на 25 декабря
	events := cal.GetEventsForDay(userID, date1)
	if len(events) != 1 {
		t.Fatalf("expected 1 event for day, got %d", len(events))
	}
	if events[0].Title != "Christmas" {
		t.Fatalf("expected event title to be 'Christmas', got %s", events[0].Title)
	}
}

func TestGetEventsForWeek(t *testing.T) {
	cal := NewCalendar()
	userID := 1

	// Создаем события на разные дни одной недели
	date1 := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC) // Понедельник
	date2 := time.Date(2023, 12, 26, 0, 0, 0, 0, time.UTC) // Вторник
	date3 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)   // Другая неделя

	cal.CreateEvent(userID, date1, "Christmas")
	cal.CreateEvent(userID, date2, "Boxing Day")
	cal.CreateEvent(userID, date3, "New Year")

	// Получаем события на неделю 25 декабря
	events := cal.GetEventsForWeek(userID, date1)
	if len(events) != 2 {
		t.Fatalf("expected 2 events for week, got %d", len(events))
	}
}

func TestGetEventsForMonth(t *testing.T) {
	cal := NewCalendar()
	userID := 1

	// Создаем события на разные месяцы
	date1 := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2023, 12, 26, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	cal.CreateEvent(userID, date1, "Christmas")
	cal.CreateEvent(userID, date2, "Boxing Day")
	cal.CreateEvent(userID, date3, "New Year")

	// Получаем события на декабрь 2023
	events := cal.GetEventsForMonth(userID, date1)
	if len(events) != 2 {
		t.Fatalf("expected 2 events for month, got %d", len(events))
	}
}

func TestIsSameDay(t *testing.T) {
	date1 := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	date2 := time.Date(2023, 12, 25, 15, 45, 0, 0, time.UTC)
	date3 := time.Date(2023, 12, 26, 10, 30, 0, 0, time.UTC)

	if !isSameDay(date1, date2) {
		t.Fatal("expected dates to be the same day")
	}

	if isSameDay(date1, date3) {
		t.Fatal("expected dates to be different days")
	}
}

func TestCalendarConcurrency(t *testing.T) {
	cal := NewCalendar()
	userID := 1
	date := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	title := "Christmas"

	// Запускаем несколько горутин для тестирования thread-safety
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			cal.CreateEvent(userID+id, date, title)
			done <- true
		}(i)
	}

	// Ждем завершения всех горутин
	for i := 0; i < 10; i++ {
		<-done
	}

	// Проверяем, что все события созданы
	if len(cal.events) != 10 {
		t.Fatalf("expected 10 events, got %d", len(cal.events))
	}
}
