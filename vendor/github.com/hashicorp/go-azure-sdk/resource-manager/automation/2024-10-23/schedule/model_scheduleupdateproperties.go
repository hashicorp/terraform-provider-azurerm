package schedule

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleUpdateProperties struct {
	Description *string `json:"description,omitempty"`
	IsEnabled   *bool   `json:"isEnabled,omitempty"`
}
