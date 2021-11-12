package hybridconnections

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type HybridConnectionProperties struct {
	CreatedAt                   *string `json:"createdAt,omitempty"`
	ListenerCount               *int64  `json:"listenerCount,omitempty"`
	RequiresClientAuthorization *bool   `json:"requiresClientAuthorization,omitempty"`
	UpdatedAt                   *string `json:"updatedAt,omitempty"`
	UserMetadata                *string `json:"userMetadata,omitempty"`
}

func (o HybridConnectionProperties) GetCreatedAtAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o HybridConnectionProperties) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o HybridConnectionProperties) GetUpdatedAtAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(o.UpdatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o HybridConnectionProperties) SetUpdatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedAt = &formatted
}
