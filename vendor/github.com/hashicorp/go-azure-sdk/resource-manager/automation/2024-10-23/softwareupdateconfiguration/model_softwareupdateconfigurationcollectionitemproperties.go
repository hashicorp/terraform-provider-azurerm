package softwareupdateconfiguration

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationCollectionItemProperties struct {
	CreationTime        *string                           `json:"creationTime,omitempty"`
	Frequency           *ScheduleFrequency                `json:"frequency,omitempty"`
	LastModifiedTime    *string                           `json:"lastModifiedTime,omitempty"`
	NextRun             *string                           `json:"nextRun,omitempty"`
	ProvisioningState   *string                           `json:"provisioningState,omitempty"`
	StartTime           *string                           `json:"startTime,omitempty"`
	Tasks               *SoftwareUpdateConfigurationTasks `json:"tasks,omitempty"`
	UpdateConfiguration *UpdateConfiguration              `json:"updateConfiguration,omitempty"`
}

func (o *SoftwareUpdateConfigurationCollectionItemProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SoftwareUpdateConfigurationCollectionItemProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *SoftwareUpdateConfigurationCollectionItemProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SoftwareUpdateConfigurationCollectionItemProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}

func (o *SoftwareUpdateConfigurationCollectionItemProperties) GetNextRunAsTime() (*time.Time, error) {
	if o.NextRun == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NextRun, "2006-01-02T15:04:05Z07:00")
}

func (o *SoftwareUpdateConfigurationCollectionItemProperties) SetNextRunAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NextRun = &formatted
}

func (o *SoftwareUpdateConfigurationCollectionItemProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SoftwareUpdateConfigurationCollectionItemProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
