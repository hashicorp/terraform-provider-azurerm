package consumergroups

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/formatting"
)

type ConsumerGroupProperties struct {
	CreatedAt    *string `json:"createdAt,omitempty"`
	UpdatedAt    *string `json:"updatedAt,omitempty"`
	UserMetadata *string `json:"userMetadata,omitempty"`
}

func (o ConsumerGroupProperties) ListCreatedAtAsTime() (*time.Time, error) {
	return formatting.ParseAsDateFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o ConsumerGroupProperties) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o ConsumerGroupProperties) ListUpdatedAtAsTime() (*time.Time, error) {
	return formatting.ParseAsDateFormat(o.UpdatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o ConsumerGroupProperties) SetUpdatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedAt = &formatted
}
