package fluxconfiguration

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ObjectStatusConditionDefinition struct {
	LastTransitionTime *string `json:"lastTransitionTime,omitempty"`
	Message            *string `json:"message,omitempty"`
	Reason             *string `json:"reason,omitempty"`
	Status             *string `json:"status,omitempty"`
	Type               *string `json:"type,omitempty"`
}

func (o *ObjectStatusConditionDefinition) GetLastTransitionTimeAsTime() (*time.Time, error) {
	if o.LastTransitionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastTransitionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ObjectStatusConditionDefinition) SetLastTransitionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastTransitionTime = &formatted
}
