package amlfilesystems

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AmlFilesystemArchiveStatus struct {
	ErrorCode          *string            `json:"errorCode,omitempty"`
	ErrorMessage       *string            `json:"errorMessage,omitempty"`
	LastCompletionTime *string            `json:"lastCompletionTime,omitempty"`
	LastStartedTime    *string            `json:"lastStartedTime,omitempty"`
	PercentComplete    *int64             `json:"percentComplete,omitempty"`
	State              *ArchiveStatusType `json:"state,omitempty"`
}

func (o *AmlFilesystemArchiveStatus) GetLastCompletionTimeAsTime() (*time.Time, error) {
	if o.LastCompletionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastCompletionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AmlFilesystemArchiveStatus) SetLastCompletionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastCompletionTime = &formatted
}

func (o *AmlFilesystemArchiveStatus) GetLastStartedTimeAsTime() (*time.Time, error) {
	if o.LastStartedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStartedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AmlFilesystemArchiveStatus) SetLastStartedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStartedTime = &formatted
}
