package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoscaleSchedule struct {
	Days            *[]DaysOfWeek             `json:"days,omitempty"`
	TimeAndCapacity *AutoscaleTimeAndCapacity `json:"timeAndCapacity,omitempty"`
}
