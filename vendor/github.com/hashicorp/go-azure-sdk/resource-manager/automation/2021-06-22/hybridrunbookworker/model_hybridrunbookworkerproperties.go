package hybridrunbookworker

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HybridRunbookWorkerProperties struct {
	IP                 *string     `json:"ip,omitempty"`
	LastSeenDateTime   *string     `json:"lastSeenDateTime,omitempty"`
	RegisteredDateTime *string     `json:"registeredDateTime,omitempty"`
	VMResourceId       *string     `json:"vmResourceId,omitempty"`
	WorkerName         *string     `json:"workerName,omitempty"`
	WorkerType         *WorkerType `json:"workerType,omitempty"`
}

func (o *HybridRunbookWorkerProperties) GetLastSeenDateTimeAsTime() (*time.Time, error) {
	if o.LastSeenDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSeenDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *HybridRunbookWorkerProperties) SetLastSeenDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSeenDateTime = &formatted
}

func (o *HybridRunbookWorkerProperties) GetRegisteredDateTimeAsTime() (*time.Time, error) {
	if o.RegisteredDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RegisteredDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *HybridRunbookWorkerProperties) SetRegisteredDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RegisteredDateTime = &formatted
}
