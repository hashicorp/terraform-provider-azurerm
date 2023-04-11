package jobschedule

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobScheduleCreateProperties struct {
	Parameters *map[string]string          `json:"parameters,omitempty"`
	RunOn      *string                     `json:"runOn,omitempty"`
	Runbook    RunbookAssociationProperty  `json:"runbook"`
	Schedule   ScheduleAssociationProperty `json:"schedule"`
}
