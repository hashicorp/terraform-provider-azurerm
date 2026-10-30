package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinuxParameters struct {
	ClassificationsToInclude  *[]VMGuestPatchClassificationLinux `json:"classificationsToInclude,omitempty"`
	MaintenanceRunId          *string                            `json:"maintenanceRunId,omitempty"`
	PackageNameMasksToExclude *[]string                          `json:"packageNameMasksToExclude,omitempty"`
	PackageNameMasksToInclude *[]string                          `json:"packageNameMasksToInclude,omitempty"`
}
