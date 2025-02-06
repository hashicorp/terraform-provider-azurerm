package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomaticRepairsPolicy struct {
	Enabled      *bool         `json:"enabled,omitempty"`
	GracePeriod  *string       `json:"gracePeriod,omitempty"`
	RepairAction *RepairAction `json:"repairAction,omitempty"`
}
