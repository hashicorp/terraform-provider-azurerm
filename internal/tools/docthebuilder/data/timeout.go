package data

import (
	"fmt"
	"strings"
)

type TimeoutType string

const (
	TimeoutTypeCreate TimeoutType = "create"
	TimeoutTypeRead   TimeoutType = "read"
	TimeoutTypeUpdate TimeoutType = "update"
	TimeoutTypeDelete TimeoutType = "delete"
)

var documentationVerb = map[TimeoutType]string{
	TimeoutTypeCreate: "creating",
	TimeoutTypeRead:   "retrieving",
	TimeoutTypeUpdate: "updating",
	TimeoutTypeDelete: "deleting",
}

type Timeout struct {
	Type     TimeoutType
	Duration int
	Name     string
}

func (t Timeout) String() string {
	return fmt.Sprintf("* `%s` - (Defaults to %s) Used when %s the %s", strings.ToLower(string(t.Type)), t.duration(), documentationVerb[t.Type], t.Name)
}

func (t Timeout) duration() string {
	duration := t.Duration
	unit := "minute"

	if (duration % 60) == 0 {
		duration /= 60
		unit = "hour"
	}

	if duration > 1 {
		unit += "s"
	}

	return fmt.Sprintf("%d %s", duration, unit)
}
