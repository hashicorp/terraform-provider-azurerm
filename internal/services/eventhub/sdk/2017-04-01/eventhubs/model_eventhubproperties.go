package eventhubs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/formatting"
)

type EventhubProperties struct {
	CaptureDescription     *CaptureDescription `json:"captureDescription,omitempty"`
	CreatedAt              *string             `json:"createdAt,omitempty"`
	MessageRetentionInDays *int64              `json:"messageRetentionInDays,omitempty"`
	PartitionCount         *int64              `json:"partitionCount,omitempty"`
	PartitionIds           *[]string           `json:"partitionIds,omitempty"`
	Status                 *EntityStatus       `json:"status,omitempty"`
	UpdatedAt              *string             `json:"updatedAt,omitempty"`
}

func (o EventhubProperties) ListCreatedAtAsTime() (*time.Time, error) {
	return formatting.ParseAsDateFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o EventhubProperties) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o EventhubProperties) ListUpdatedAtAsTime() (*time.Time, error) {
	return formatting.ParseAsDateFormat(o.UpdatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o EventhubProperties) SetUpdatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedAt = &formatted
}
