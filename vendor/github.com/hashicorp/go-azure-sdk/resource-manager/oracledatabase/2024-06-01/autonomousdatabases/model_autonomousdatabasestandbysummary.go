package autonomousdatabases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDatabaseStandbySummary struct {
	LagTimeInSeconds                *int64                            `json:"lagTimeInSeconds,omitempty"`
	LifecycleDetails                *string                           `json:"lifecycleDetails,omitempty"`
	LifecycleState                  *AutonomousDatabaseLifecycleState `json:"lifecycleState,omitempty"`
	TimeDataGuardRoleChanged        *string                           `json:"timeDataGuardRoleChanged,omitempty"`
	TimeDisasterRecoveryRoleChanged *string                           `json:"timeDisasterRecoveryRoleChanged,omitempty"`
}
