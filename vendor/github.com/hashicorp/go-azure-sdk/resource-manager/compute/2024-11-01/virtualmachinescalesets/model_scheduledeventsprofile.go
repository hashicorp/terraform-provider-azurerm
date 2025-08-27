package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduledEventsProfile struct {
	OsImageNotificationProfile   *OSImageNotificationProfile   `json:"osImageNotificationProfile,omitempty"`
	TerminateNotificationProfile *TerminateNotificationProfile `json:"terminateNotificationProfile,omitempty"`
}
