package updatesummaries

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrecheckResult struct {
	AdditionalData     *string             `json:"additionalData,omitempty"`
	Description        *string             `json:"description,omitempty"`
	HealthCheckSource  *string             `json:"healthCheckSource,omitempty"`
	Name               *string             `json:"name,omitempty"`
	Remediation        *string             `json:"remediation,omitempty"`
	Severity           *Severity           `json:"severity,omitempty"`
	Status             *Status             `json:"status,omitempty"`
	Tags               *PrecheckResultTags `json:"tags,omitempty"`
	TargetResourceID   *string             `json:"targetResourceID,omitempty"`
	TargetResourceName *string             `json:"targetResourceName,omitempty"`
	Timestamp          *string             `json:"timestamp,omitempty"`
	Title              *string             `json:"title,omitempty"`
}

func (o *PrecheckResult) GetTimestampAsTime() (*time.Time, error) {
	if o.Timestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *PrecheckResult) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = &formatted
}
