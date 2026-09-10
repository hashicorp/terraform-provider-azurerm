package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationUsage struct {
	JobsSummary                      *JobsSummary       `json:"jobsSummary,omitempty"`
	MonitoringSummary                *MonitoringSummary `json:"monitoringSummary,omitempty"`
	ProtectedItemCount               *int64             `json:"protectedItemCount,omitempty"`
	RecoveryPlanCount                *int64             `json:"recoveryPlanCount,omitempty"`
	RecoveryServicesProviderAuthType *int64             `json:"recoveryServicesProviderAuthType,omitempty"`
	RegisteredServersCount           *int64             `json:"registeredServersCount,omitempty"`
}
