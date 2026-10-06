package serverendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudTieringDatePolicyStatus struct {
	LastUpdatedTimestamp                 *string `json:"lastUpdatedTimestamp,omitempty"`
	TieredFilesMostRecentAccessTimestamp *string `json:"tieredFilesMostRecentAccessTimestamp,omitempty"`
}

func (o *CloudTieringDatePolicyStatus) GetLastUpdatedTimestampAsTime() (*time.Time, error) {
	if o.LastUpdatedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudTieringDatePolicyStatus) SetLastUpdatedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTimestamp = &formatted
}

func (o *CloudTieringDatePolicyStatus) GetTieredFilesMostRecentAccessTimestampAsTime() (*time.Time, error) {
	if o.TieredFilesMostRecentAccessTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TieredFilesMostRecentAccessTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudTieringDatePolicyStatus) SetTieredFilesMostRecentAccessTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TieredFilesMostRecentAccessTimestamp = &formatted
}
