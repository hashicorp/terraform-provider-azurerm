package cluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Notification struct {
	IsEnabled            bool                 `json:"isEnabled"`
	NotificationCategory NotificationCategory `json:"notificationCategory"`
	NotificationLevel    NotificationLevel    `json:"notificationLevel"`
	NotificationTargets  []NotificationTarget `json:"notificationTargets"`
}
