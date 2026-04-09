package serverendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEndpointSyncSessionStatus struct {
	FilesNotSyncingErrors          *[]ServerEndpointFilesNotSyncingError `json:"filesNotSyncingErrors,omitempty"`
	LastSyncMode                   *ServerEndpointSyncMode               `json:"lastSyncMode,omitempty"`
	LastSyncPerItemErrorCount      *int64                                `json:"lastSyncPerItemErrorCount,omitempty"`
	LastSyncResult                 *int64                                `json:"lastSyncResult,omitempty"`
	LastSyncSuccessTimestamp       *string                               `json:"lastSyncSuccessTimestamp,omitempty"`
	LastSyncTimestamp              *string                               `json:"lastSyncTimestamp,omitempty"`
	PersistentFilesNotSyncingCount *int64                                `json:"persistentFilesNotSyncingCount,omitempty"`
	TransientFilesNotSyncingCount  *int64                                `json:"transientFilesNotSyncingCount,omitempty"`
}

func (o *ServerEndpointSyncSessionStatus) GetLastSyncSuccessTimestampAsTime() (*time.Time, error) {
	if o.LastSyncSuccessTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSyncSuccessTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerEndpointSyncSessionStatus) SetLastSyncSuccessTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSyncSuccessTimestamp = &formatted
}

func (o *ServerEndpointSyncSessionStatus) GetLastSyncTimestampAsTime() (*time.Time, error) {
	if o.LastSyncTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSyncTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerEndpointSyncSessionStatus) SetLastSyncTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSyncTimestamp = &formatted
}
