package managedcassandras

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedCassandraReaperStatus struct {
	Healthy         *bool              `json:"healthy,omitempty"`
	RepairRunIds    *map[string]string `json:"repairRunIds,omitempty"`
	RepairSchedules *map[string]string `json:"repairSchedules,omitempty"`
}
