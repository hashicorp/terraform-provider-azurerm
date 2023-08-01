package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeStartStopSchedule struct {
	Action             *ComputePowerAction `json:"action,omitempty"`
	Cron               *Cron               `json:"cron,omitempty"`
	Id                 *string             `json:"id,omitempty"`
	ProvisioningStatus *ProvisioningStatus `json:"provisioningStatus,omitempty"`
	Recurrence         *Recurrence         `json:"recurrence,omitempty"`
	Schedule           *ScheduleBase       `json:"schedule,omitempty"`
	Status             *ScheduleStatus     `json:"status,omitempty"`
	TriggerType        *TriggerType        `json:"triggerType,omitempty"`
}
