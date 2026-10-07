package dbservers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DbServerPatchingDetails struct {
	EstimatedPatchDuration *int64                  `json:"estimatedPatchDuration,omitempty"`
	PatchingStatus         *DbServerPatchingStatus `json:"patchingStatus,omitempty"`
	TimePatchingEnded      *string                 `json:"timePatchingEnded,omitempty"`
	TimePatchingStarted    *string                 `json:"timePatchingStarted,omitempty"`
}

func (o *DbServerPatchingDetails) GetTimePatchingEndedAsTime() (*time.Time, error) {
	if o.TimePatchingEnded == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimePatchingEnded, "2006-01-02T15:04:05Z07:00")
}

func (o *DbServerPatchingDetails) SetTimePatchingEndedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimePatchingEnded = &formatted
}

func (o *DbServerPatchingDetails) GetTimePatchingStartedAsTime() (*time.Time, error) {
	if o.TimePatchingStarted == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimePatchingStarted, "2006-01-02T15:04:05Z07:00")
}

func (o *DbServerPatchingDetails) SetTimePatchingStartedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimePatchingStarted = &formatted
}
