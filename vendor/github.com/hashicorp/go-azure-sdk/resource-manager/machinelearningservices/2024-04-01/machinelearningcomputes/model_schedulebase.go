package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleBase struct {
	Id                 *string                    `json:"id,omitempty"`
	ProvisioningStatus *ScheduleProvisioningState `json:"provisioningStatus,omitempty"`
	Status             *ScheduleStatus            `json:"status,omitempty"`
}
