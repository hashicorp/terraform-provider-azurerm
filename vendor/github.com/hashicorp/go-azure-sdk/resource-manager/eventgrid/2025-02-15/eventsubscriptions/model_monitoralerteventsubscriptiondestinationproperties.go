package eventsubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorAlertEventSubscriptionDestinationProperties struct {
	ActionGroups *[]string             `json:"actionGroups,omitempty"`
	Description  *string               `json:"description,omitempty"`
	Severity     *MonitorAlertSeverity `json:"severity,omitempty"`
}
