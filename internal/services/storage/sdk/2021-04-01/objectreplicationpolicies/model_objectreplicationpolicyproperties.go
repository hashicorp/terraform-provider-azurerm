package objectreplicationpolicies

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type ObjectReplicationPolicyProperties struct {
	DestinationAccount string                         `json:"destinationAccount"`
	EnabledTime        *string                        `json:"enabledTime,omitempty"`
	PolicyId           *string                        `json:"policyId,omitempty"`
	Rules              *[]ObjectReplicationPolicyRule `json:"rules,omitempty"`
	SourceAccount      string                         `json:"sourceAccount"`
}

func (o ObjectReplicationPolicyProperties) GetEnabledTimeAsTime() (*time.Time, error) {
	if o.EnabledTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EnabledTime, "2006-01-02T15:04:05Z07:00")
}

func (o ObjectReplicationPolicyProperties) SetEnabledTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EnabledTime = &formatted
}
