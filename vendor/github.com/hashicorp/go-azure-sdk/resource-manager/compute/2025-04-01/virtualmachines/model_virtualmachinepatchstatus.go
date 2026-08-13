package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachinePatchStatus struct {
	AvailablePatchSummary        *AvailablePatchSummary        `json:"availablePatchSummary,omitempty"`
	ConfigurationStatuses        *[]InstanceViewStatus         `json:"configurationStatuses,omitempty"`
	LastPatchInstallationSummary *LastPatchInstallationSummary `json:"lastPatchInstallationSummary,omitempty"`
}
