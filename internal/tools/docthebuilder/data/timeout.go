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

var friendlyString = map[TimeoutType]string{
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
	return fmt.Sprintf("* `%s` - (Defaults to %s) Used when %s the %s", strings.ToLower(string(t.Type)), t.duration(), friendlyString[t.Type], t.Name)
}

func (t Timeout) duration() string {
	duration := t.Duration
	unit := "minute"

	if (duration % 60) == 0 {
		duration = duration / 60
		unit = "hour"
	}

	if duration > 1 {
		unit = unit + "s"
	}

	return fmt.Sprintf("%d %s", duration, unit)
}

func getTimeouts(rd *ResourceData) {
	t := rd.Resource.Timeouts
	if t == nil {
		return
	}

	if t.Create != nil {
		rd.Timeouts = append(rd.Timeouts, Timeout{
			Type:     TimeoutTypeCreate,
			Duration: int(t.Create.Minutes()),
			Name:     "TODO",
		})
	}

	if t.Read != nil {
		rd.Timeouts = append(rd.Timeouts, Timeout{
			Type:     TimeoutTypeRead,
			Duration: int(t.Read.Minutes()),
			Name:     "TODO",
		})
	}

	if t.Update != nil {
		rd.Timeouts = append(rd.Timeouts, Timeout{
			Type:     TimeoutTypeUpdate,
			Duration: int(t.Update.Minutes()),
			Name:     "TODO",
		})
	}

	if t.Delete != nil {
		rd.Timeouts = append(rd.Timeouts, Timeout{
			Type:     TimeoutTypeDelete,
			Duration: int(t.Delete.Minutes()),
			Name:     "TODO",
		})
	}
}
