package events

import (
	"errors"
	"regexp"
	"time"

	"github.com/araddon/dateparse"
)

func IsValidTitle(title string) bool {
	pattern := "^[a-zA-Z0-9а-яА-Я ,/.]{3,55}$"
	matched, err := regexp.MatchString(pattern, title)
	if err != nil {
		return false
	}
	return matched
}
func IsValidDate(dateStr string) (time.Time, error) {
	t, err := dateparse.ParseLocal(dateStr)
	if err != nil {
		return t, errors.New("неверный формат даты")
	}
	return t, nil
}
