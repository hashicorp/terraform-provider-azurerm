package sourcecontrolsyncjob

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceControlSyncJobByIdProperties struct {
	CreationTime           *string            `json:"creationTime,omitempty"`
	EndTime                *string            `json:"endTime,omitempty"`
	Exception              *string            `json:"exception,omitempty"`
	ProvisioningState      *ProvisioningState `json:"provisioningState,omitempty"`
	SourceControlSyncJobId *string            `json:"sourceControlSyncJobId,omitempty"`
	StartTime              *string            `json:"startTime,omitempty"`
	SyncType               *SyncType          `json:"syncType,omitempty"`
}

func (o *SourceControlSyncJobByIdProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SourceControlSyncJobByIdProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *SourceControlSyncJobByIdProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SourceControlSyncJobByIdProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *SourceControlSyncJobByIdProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SourceControlSyncJobByIdProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
