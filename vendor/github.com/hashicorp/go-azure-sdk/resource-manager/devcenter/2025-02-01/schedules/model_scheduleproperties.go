package schedules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleProperties struct {
	Frequency         *ScheduledFrequency   `json:"frequency,omitempty"`
	Location          *string               `json:"location,omitempty"`
	ProvisioningState *ProvisioningState    `json:"provisioningState,omitempty"`
	State             *ScheduleEnableStatus `json:"state,omitempty"`
	Tags              *map[string]string    `json:"tags,omitempty"`
	Time              *string               `json:"time,omitempty"`
	TimeZone          *string               `json:"timeZone,omitempty"`
	Type              *ScheduledType        `json:"type,omitempty"`
}
