package scheduledactions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduledActionProperties struct {
	DisplayName     string                 `json:"displayName"`
	FileDestination *FileDestination       `json:"fileDestination,omitempty"`
	Notification    NotificationProperties `json:"notification"`
	Schedule        ScheduleProperties     `json:"schedule"`
	Scope           *string                `json:"scope,omitempty"`
	Status          ScheduledActionStatus  `json:"status"`
	ViewId          string                 `json:"viewId"`
}
