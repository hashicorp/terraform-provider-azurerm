package schedules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleProperties struct {
	Frequency         ScheduledFrequency    `json:"frequency"`
	ProvisioningState *ProvisioningState    `json:"provisioningState,omitempty"`
	State             *ScheduleEnableStatus `json:"state,omitempty"`
	Time              string                `json:"time"`
	TimeZone          string                `json:"timeZone"`
	Type              ScheduledType         `json:"type"`
}
