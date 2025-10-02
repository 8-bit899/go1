package calendar

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/io893/calendar_app/events"
	"github.com/io893/calendar_app/storage"
)

type Calendar struct {
	calendarEvents map[string]*events.Event
	storage        storage.Store
	Notification   chan string `json:"-"`
}

func (c *Calendar) Save() error {
	data, err := json.Marshal(c.calendarEvents)
	if err != nil {

		return fmt.Errorf("error marshal %w", err)
	}
	err = c.storage.Save(data)
	return err
}
func (c *Calendar) Load() error {
	data, err := c.storage.Load()
	if err != nil {
	}

	err = json.Unmarshal(data, &c.calendarEvents)
	if err != nil {

		return err
	}
	return err
}
func NewCalendar(s storage.Store) *Calendar {
	return &Calendar{
		calendarEvents: make(map[string]*events.Event),
		storage:        s,
		Notification:   make(chan string),
	}
}
func (c *Calendar) AddEvent(title string, date string, priority events.Priority) (events.Event, error) {

	e, err := events.NewEvent(title, date, priority)
	if err != nil {
		return *e, err
	}
	c.calendarEvents[e.ID] = e

	return *e, nil
}
func stringOfThree(key string, title string, date string, priority string) string {
	return key + ": " + title + " >> " + priority + " >> " + date

}
func (c *Calendar) ShowEvents() {
	for _, e := range c.calendarEvents {
		c.Notify(stringOfThree(e.ID, e.Title, e.StartAt.String(), string(e.Priority)))
	}
}
func (c *Calendar) EditEvent(key string, title string, date string, priority events.Priority) error {
	_, ok := c.calendarEvents[key]
	if ok {
		err := c.calendarEvents[key].UpdateEvent(title, date, priority)
		if err != nil {
			return fmt.Errorf(events.ErrorUpdatingEventMessage+" %w", err)

		}
		return nil
	} else {
		return errors.New(events.EventNotFoundMessage)
	}
}
func (c *Calendar) DeleteEvent(key string) error {
	_, ok := c.calendarEvents[key]
	if ok {
		delete(c.calendarEvents, key)
		return nil //events.EventDeletedMessage)
	} else {
		return errors.New(events.EventNotFoundMessage)
	}
}
func (c *Calendar) SetEventReminder(key string, msg string, at string) error {
	_, ok := c.calendarEvents[key]
	if ok {
		err := c.calendarEvents[key].AddReminder(msg, at, c.Notify)
		if err != nil {
			return err
		}

		return nil
	}
	return fmt.Errorf(events.ReminderNotAdd+"%w", errors.New(events.EventNotFoundMessage))

}
func (c *Calendar) CancelEventReminder(key string) error {
	_, ok := c.calendarEvents[key]
	if ok {
		err := c.calendarEvents[key].RemoveReminder()
		if err != nil {
			return err
		}
		c.Notify(events.ReminderRemov)
		return nil
	}
	return fmt.Errorf(events.ReminderCannotDeleted+"%w", errors.New(events.EventNotFoundMessage))
}
func (c *Calendar) Notify(msg string) {
	c.Notification <- msg
}
func (c *Calendar) CloseNotify() {
	close(c.Notification)
}
