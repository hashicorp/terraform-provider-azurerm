package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureMonitorAlertSettings struct {
	AlertsForAllFailoverIssues    *AlertsState `json:"alertsForAllFailoverIssues,omitempty"`
	AlertsForAllJobFailures       *AlertsState `json:"alertsForAllJobFailures,omitempty"`
	AlertsForAllReplicationIssues *AlertsState `json:"alertsForAllReplicationIssues,omitempty"`
}
