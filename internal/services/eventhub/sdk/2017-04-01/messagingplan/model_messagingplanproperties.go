package messagingplan

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type MessagingPlanProperties struct {
	Revision             *int64  `json:"revision,omitempty"`
	SelectedEventHubUnit *int64  `json:"selectedEventHubUnit,omitempty"`
	Sku                  *int64  `json:"sku,omitempty"`
	UpdatedAt            *string `json:"updatedAt,omitempty"`
}

func (o MessagingPlanProperties) GetUpdatedAtAsTime() (*time.Time, error) {
	if o.UpdatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o MessagingPlanProperties) SetUpdatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedAt = &formatted
}
