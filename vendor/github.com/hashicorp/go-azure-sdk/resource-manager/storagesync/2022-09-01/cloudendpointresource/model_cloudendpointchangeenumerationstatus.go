package cloudendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudEndpointChangeEnumerationStatus struct {
	Activity              *CloudEndpointChangeEnumerationActivity   `json:"activity,omitempty"`
	LastEnumerationStatus *CloudEndpointLastChangeEnumerationStatus `json:"lastEnumerationStatus,omitempty"`
	LastUpdatedTimestamp  *string                                   `json:"lastUpdatedTimestamp,omitempty"`
}

func (o *CloudEndpointChangeEnumerationStatus) GetLastUpdatedTimestampAsTime() (*time.Time, error) {
	if o.LastUpdatedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudEndpointChangeEnumerationStatus) SetLastUpdatedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTimestamp = &formatted
}
