package autonomousdatabaseversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDbVersionProperties struct {
	DbWorkload        *WorkloadType `json:"dbWorkload,omitempty"`
	IsDefaultForFree  *bool         `json:"isDefaultForFree,omitempty"`
	IsDefaultForPaid  *bool         `json:"isDefaultForPaid,omitempty"`
	IsFreeTierEnabled *bool         `json:"isFreeTierEnabled,omitempty"`
	IsPaidEnabled     *bool         `json:"isPaidEnabled,omitempty"`
	Version           string        `json:"version"`
}
