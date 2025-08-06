package serverendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEndpointRecallStatus struct {
	LastUpdatedTimestamp   *string                      `json:"lastUpdatedTimestamp,omitempty"`
	RecallErrors           *[]ServerEndpointRecallError `json:"recallErrors,omitempty"`
	TotalRecallErrorsCount *int64                       `json:"totalRecallErrorsCount,omitempty"`
}

func (o *ServerEndpointRecallStatus) GetLastUpdatedTimestampAsTime() (*time.Time, error) {
	if o.LastUpdatedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerEndpointRecallStatus) SetLastUpdatedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTimestamp = &formatted
}
