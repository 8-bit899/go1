package reminder

import (
	"fmt"
	"time"
)

type Reminder struct {
	Message string
	At      time.Time
	Sent    bool
	Timer   *time.Timer  `json:"-"`
	Notify  func(string) `json:"-"`
}

func NewReminder(message string, at time.Time, notify func(string)) *Reminder {

	return &Reminder{
		Message: message,
		At:      at,
		Sent:    false,
		Notify:  notify,
		Timer:   nil,
	}
}
func (r *Reminder) Send() {
	if r.Sent {
		fmt.Println("Напоминание уже отправлено")
		return
	}
	r.Notify("Reminder: " + r.Message)
	r.Sent = true
}
func (r *Reminder) Start() {
	duration := time.Until(r.At)
	if duration <= 0 {
		fmt.Println("Некорректное время уведомления")
		return
	}
	r.Timer = time.AfterFunc(duration, r.Send)

}
func (r *Reminder) Stop() {
	if r.Sent && r.Timer != nil {
		r.Timer.Stop()
	}

}
