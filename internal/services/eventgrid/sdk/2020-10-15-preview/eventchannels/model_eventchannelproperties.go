package eventchannels

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type EventChannelProperties struct {
	Destination                     *EventChannelDestination       `json:"destination,omitempty"`
	ExpirationTimeIfNotActivatedUtc *string                        `json:"expirationTimeIfNotActivatedUtc,omitempty"`
	Filter                          *EventChannelFilter            `json:"filter,omitempty"`
	PartnerTopicFriendlyDescription *string                        `json:"partnerTopicFriendlyDescription,omitempty"`
	PartnerTopicReadinessState      *PartnerTopicReadinessState    `json:"partnerTopicReadinessState,omitempty"`
	ProvisioningState               *EventChannelProvisioningState `json:"provisioningState,omitempty"`
	Source                          *EventChannelSource            `json:"source,omitempty"`
}

func (o EventChannelProperties) GetExpirationTimeIfNotActivatedUtcAsTime() (*time.Time, error) {
	if o.ExpirationTimeIfNotActivatedUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationTimeIfNotActivatedUtc, "2006-01-02T15:04:05Z07:00")
}

func (o EventChannelProperties) SetExpirationTimeIfNotActivatedUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationTimeIfNotActivatedUtc = &formatted
}
