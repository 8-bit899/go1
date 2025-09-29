package events

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/io893/calendar_app/reminder"
)

type Event struct {
	ID       string
	Title    string
	StartAt  time.Time
	Priority Priority
	Reminder *reminder.Reminder
}

func getNextID() string {
	return uuid.New().String()
}
func NewEvent(title string, dateStr string, priority Priority) (*Event, error) {
	t, err := IsValidDate(dateStr)
	if err != nil {
		return &Event{}, err
	}
	if !IsValidTitle(title) {
		return &Event{}, errors.New(InvalidHeaderFormatMessage)
	}
	err = priority.Validate()
	if err != nil {
		return &Event{}, err
	}
	return &Event{
		ID:       getNextID(),
		Title:    title,
		StartAt:  t,
		Priority: priority,
		Reminder: nil,
	}, nil
}
func (e *Event) UpdateEvent(title string, data string, priority Priority) error {
	if !IsValidTitle(title) {
		return errors.New(InvalidHeaderFormatMessage)

	}
	t, err := IsValidDate(data)
	if err != nil {
		return err
	}
	err = priority.Validate()
	if err != nil {
		return err
	}
	e.Priority = priority
	e.StartAt = t
	e.Title = title
	return nil
}
func (e *Event) AddReminder(message string, data string, notify func(string)) {
	at, err := IsValidDate(data)
	if err != nil {
		fmt.Println("проблемы с датой в напоминании")
	}
	e.Reminder = reminder.NewReminder(message, at, notify)
	e.Reminder.Start()
}
func (e *Event) RemoveReminder() {
	if !e.Reminder.Sent {
		fmt.Println("нельзя удалить то чего нет. уведомление в событии отсутствует. ")
	}
	e.Reminder.Stop()
	e.Reminder = nil
	fmt.Println("уведомление удалено")
}
