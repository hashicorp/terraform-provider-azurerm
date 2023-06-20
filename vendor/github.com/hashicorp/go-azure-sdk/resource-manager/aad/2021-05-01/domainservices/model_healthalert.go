package domainservices

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HealthAlert struct {
	Id            *string `json:"id,omitempty"`
	Issue         *string `json:"issue,omitempty"`
	LastDetected  *string `json:"lastDetected,omitempty"`
	Name          *string `json:"name,omitempty"`
	Raised        *string `json:"raised,omitempty"`
	ResolutionUri *string `json:"resolutionUri,omitempty"`
	Severity      *string `json:"severity,omitempty"`
}

func (o *HealthAlert) GetLastDetectedAsTime() (*time.Time, error) {
	if o.LastDetected == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastDetected, "2006-01-02T15:04:05Z07:00")
}

func (o *HealthAlert) SetLastDetectedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastDetected = &formatted
}

func (o *HealthAlert) GetRaisedAsTime() (*time.Time, error) {
	if o.Raised == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Raised, "2006-01-02T15:04:05Z07:00")
}

func (o *HealthAlert) SetRaisedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Raised = &formatted
}
