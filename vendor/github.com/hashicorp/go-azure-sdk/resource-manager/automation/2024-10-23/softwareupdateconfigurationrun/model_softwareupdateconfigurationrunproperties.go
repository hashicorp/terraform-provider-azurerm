package softwareupdateconfigurationrun

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationRunProperties struct {
	ComputerCount               *int64                               `json:"computerCount,omitempty"`
	ConfiguredDuration          *string                              `json:"configuredDuration,omitempty"`
	CreatedBy                   *string                              `json:"createdBy,omitempty"`
	CreationTime                *string                              `json:"creationTime,omitempty"`
	EndTime                     *string                              `json:"endTime,omitempty"`
	FailedCount                 *int64                               `json:"failedCount,omitempty"`
	LastModifiedBy              *string                              `json:"lastModifiedBy,omitempty"`
	LastModifiedTime            *string                              `json:"lastModifiedTime,omitempty"`
	OsType                      *string                              `json:"osType,omitempty"`
	SoftwareUpdateConfiguration *UpdateConfigurationNavigation       `json:"softwareUpdateConfiguration,omitempty"`
	StartTime                   *string                              `json:"startTime,omitempty"`
	Status                      *string                              `json:"status,omitempty"`
	Tasks                       *SoftwareUpdateConfigurationRunTasks `json:"tasks,omitempty"`
}

func (o *SoftwareUpdateConfigurationRunProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SoftwareUpdateConfigurationRunProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *SoftwareUpdateConfigurationRunProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SoftwareUpdateConfigurationRunProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *SoftwareUpdateConfigurationRunProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SoftwareUpdateConfigurationRunProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}

func (o *SoftwareUpdateConfigurationRunProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SoftwareUpdateConfigurationRunProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
