package softwareupdateconfiguration

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftwareUpdateConfigurationProperties struct {
	CreatedBy           *string                           `json:"createdBy,omitempty"`
	CreationTime        *string                           `json:"creationTime,omitempty"`
	Error               *ErrorResponse                    `json:"error,omitempty"`
	LastModifiedBy      *string                           `json:"lastModifiedBy,omitempty"`
	LastModifiedTime    *string                           `json:"lastModifiedTime,omitempty"`
	ProvisioningState   *string                           `json:"provisioningState,omitempty"`
	ScheduleInfo        SUCScheduleProperties             `json:"scheduleInfo"`
	Tasks               *SoftwareUpdateConfigurationTasks `json:"tasks,omitempty"`
	UpdateConfiguration UpdateConfiguration               `json:"updateConfiguration"`
}

func (o *SoftwareUpdateConfigurationProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SoftwareUpdateConfigurationProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *SoftwareUpdateConfigurationProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SoftwareUpdateConfigurationProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
