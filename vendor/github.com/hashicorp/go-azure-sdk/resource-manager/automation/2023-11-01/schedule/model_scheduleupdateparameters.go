package schedule

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleUpdateParameters struct {
	Name       *string                   `json:"name,omitempty"`
	Properties *ScheduleUpdateProperties `json:"properties,omitempty"`
}
