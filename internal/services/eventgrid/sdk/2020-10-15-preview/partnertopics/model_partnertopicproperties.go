package partnertopics

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type PartnerTopicProperties struct {
	ActivationState                 *PartnerTopicActivationState   `json:"activationState,omitempty"`
	ExpirationTimeIfNotActivatedUtc *string                        `json:"expirationTimeIfNotActivatedUtc,omitempty"`
	PartnerTopicFriendlyDescription *string                        `json:"partnerTopicFriendlyDescription,omitempty"`
	ProvisioningState               *PartnerTopicProvisioningState `json:"provisioningState,omitempty"`
	Source                          *string                        `json:"source,omitempty"`
}

func (o PartnerTopicProperties) GetExpirationTimeIfNotActivatedUtcAsTime() (*time.Time, error) {
	if o.ExpirationTimeIfNotActivatedUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationTimeIfNotActivatedUtc, "2006-01-02T15:04:05Z07:00")
}

func (o PartnerTopicProperties) SetExpirationTimeIfNotActivatedUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationTimeIfNotActivatedUtc = &formatted
}
