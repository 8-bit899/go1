package events

import (
	"errors"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

func (p Priority) Validate() error {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return nil
	default:
		return errors.New(InvalidPriority)
	}
}
func (p Priority) IsHigh() bool {
	switch p {
	case PriorityHigh:
		return true
	default:
		return false
	}
}
func (p Priority) IsMedium() bool {
	switch p {
	case PriorityMedium:
		return true
	default:
		return false
	}
}
func (p Priority) IsLow() bool {
	switch p {
	case PriorityLow:
		return true
	default:
		return false
	}
}
func (p Priority) All() [3]Priority {
	return [3]Priority{PriorityHigh, PriorityMedium, PriorityLow}
}

func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}
func (p Priority) Compare(pr Priority) bool {
	return p == pr
}
func (p Priority) Next() (Priority, error) {
	switch p {
	case PriorityLow:
		return PriorityMedium, nil
	case PriorityMedium:
		return PriorityHigh, nil
	case PriorityHigh:
		return PriorityLow, nil
	default:
		return p, errors.New(InvalidPriority)

	}

}
