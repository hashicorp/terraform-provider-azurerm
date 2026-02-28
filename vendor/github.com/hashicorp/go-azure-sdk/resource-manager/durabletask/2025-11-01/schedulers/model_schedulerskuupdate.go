package schedulers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchedulerSkuUpdate struct {
	Capacity        *int64            `json:"capacity,omitempty"`
	Name            *SchedulerSkuName `json:"name,omitempty"`
	RedundancyState *RedundancyState  `json:"redundancyState,omitempty"`
}
