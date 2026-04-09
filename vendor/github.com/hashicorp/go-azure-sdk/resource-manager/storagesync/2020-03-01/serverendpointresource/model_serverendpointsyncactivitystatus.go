package serverendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEndpointSyncActivityStatus struct {
	AppliedBytes      *int64                  `json:"appliedBytes,omitempty"`
	AppliedItemCount  *int64                  `json:"appliedItemCount,omitempty"`
	PerItemErrorCount *int64                  `json:"perItemErrorCount,omitempty"`
	SyncMode          *ServerEndpointSyncMode `json:"syncMode,omitempty"`
	Timestamp         *string                 `json:"timestamp,omitempty"`
	TotalBytes        *int64                  `json:"totalBytes,omitempty"`
	TotalItemCount    *int64                  `json:"totalItemCount,omitempty"`
}

func (o *ServerEndpointSyncActivityStatus) GetTimestampAsTime() (*time.Time, error) {
	if o.Timestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerEndpointSyncActivityStatus) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = &formatted
}
