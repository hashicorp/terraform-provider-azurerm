package softwareupdateconfigurationmachinerun

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateConfigurationMachineRunProperties struct {
	ConfiguredDuration          *string                        `json:"configuredDuration,omitempty"`
	CorrelationId               *string                        `json:"correlationId,omitempty"`
	CreatedBy                   *string                        `json:"createdBy,omitempty"`
	CreationTime                *string                        `json:"creationTime,omitempty"`
	EndTime                     *string                        `json:"endTime,omitempty"`
	Error                       *ErrorResponse                 `json:"error,omitempty"`
	Job                         *JobNavigation                 `json:"job,omitempty"`
	LastModifiedBy              *string                        `json:"lastModifiedBy,omitempty"`
	LastModifiedTime            *string                        `json:"lastModifiedTime,omitempty"`
	OsType                      *string                        `json:"osType,omitempty"`
	SoftwareUpdateConfiguration *UpdateConfigurationNavigation `json:"softwareUpdateConfiguration,omitempty"`
	SourceComputerId            *string                        `json:"sourceComputerId,omitempty"`
	StartTime                   *string                        `json:"startTime,omitempty"`
	Status                      *string                        `json:"status,omitempty"`
	TargetComputer              *string                        `json:"targetComputer,omitempty"`
	TargetComputerType          *string                        `json:"targetComputerType,omitempty"`
}

func (o *UpdateConfigurationMachineRunProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateConfigurationMachineRunProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *UpdateConfigurationMachineRunProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateConfigurationMachineRunProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *UpdateConfigurationMachineRunProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateConfigurationMachineRunProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}

func (o *UpdateConfigurationMachineRunProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateConfigurationMachineRunProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
