package autoimportjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoImportJobProperties struct {
	AdminStatus            *AdminStatus                   `json:"adminStatus,omitempty"`
	AutoImportPrefixes     *[]string                      `json:"autoImportPrefixes,omitempty"`
	ConflictResolutionMode *ConflictResolutionMode        `json:"conflictResolutionMode,omitempty"`
	EnableDeletions        *bool                          `json:"enableDeletions,omitempty"`
	MaximumErrors          *int64                         `json:"maximumErrors,omitempty"`
	ProvisioningState      *ProvisioningState             `json:"provisioningState,omitempty"`
	Status                 *AutoImportJobPropertiesStatus `json:"status,omitempty"`
}
