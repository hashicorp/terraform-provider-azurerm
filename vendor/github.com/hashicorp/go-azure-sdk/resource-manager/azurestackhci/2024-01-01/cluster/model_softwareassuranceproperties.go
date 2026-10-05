package cluster

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareAssuranceProperties struct {
	LastUpdated             *string                  `json:"lastUpdated,omitempty"`
	SoftwareAssuranceIntent *SoftwareAssuranceIntent `json:"softwareAssuranceIntent,omitempty"`
	SoftwareAssuranceStatus *SoftwareAssuranceStatus `json:"softwareAssuranceStatus,omitempty"`
}

func (o *SoftwareAssuranceProperties) GetLastUpdatedAsTime() (*time.Time, error) {
	if o.LastUpdated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdated, "2006-01-02T15:04:05Z07:00")
}

func (o *SoftwareAssuranceProperties) SetLastUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdated = &formatted
}
