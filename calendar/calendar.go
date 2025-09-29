package calendar

import (
	"encoding/json"

	"github.com/io893/calendar_app/events"
	"github.com/io893/calendar_app/storage"
)

type Calendar struct {
	calendarEvents map[string]*events.Event
	storage        storage.Store
	Notification   chan string
}

func (c *Calendar) Save() error {
	data, err := json.Marshal(c.calendarEvents)
	if err != nil {

		return err
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
	c.Notify(events.EventAddedMessage)
	return *e, nil
}
func stringOfThree(key string, title string, date string) string {
	return key + ": " + title + " >> " + date

}
func (c *Calendar) ShowEvents() {
	for _, e := range c.calendarEvents {
		c.Notify(stringOfThree(e.ID, e.Title, e.StartAt.String()))
	}
}
func (c *Calendar) EditEvent(key string, title string, date string, priority events.Priority) {
	_, ok := c.calendarEvents[key]
	if ok {
		err := c.calendarEvents[key].UpdateEvent(title, date, priority)
		if err != nil {
			c.Notify(events.ErrorUpdatingEventMessage)
			return
		}
		c.Notify(events.EventUpdatedMessage)
	} else {
		c.Notify(events.EventNotFoundMessage)
	}
}
func (c *Calendar) DeleteEvent(key string) {
	_, ok := c.calendarEvents[key]
	if ok {
		delete(c.calendarEvents, key)
		c.Notify(events.EventDeletedMessage)
	} else {
		c.Notify(events.EventNotFoundMessage)
	}
}
func (c *Calendar) SetEventReminder(key string, msg string, at string) {
	_, ok := c.calendarEvents[key]
	if ok {
		c.calendarEvents[key].AddReminder(msg, at, c.Notify)
		c.Notify("напоминание добавлено")
	} else {
		c.Notify(events.EventNotFoundMessage)
	}
}
func (c *Calendar) CancelEventReminder(key string) {
	_, ok := c.calendarEvents[key]
	if ok {
		c.calendarEvents[key].RemoveReminder()
	}

}
func (c *Calendar) Notify(msg string) {
	c.Notification <- msg
}
func (c *Calendar) CloseNotify() {
	close(c.Notification)
}
