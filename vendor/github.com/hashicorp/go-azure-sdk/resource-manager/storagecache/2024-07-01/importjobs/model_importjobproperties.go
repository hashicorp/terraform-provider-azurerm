package importjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImportJobProperties struct {
	AdminStatus            *ImportJobAdminStatus           `json:"adminStatus,omitempty"`
	ConflictResolutionMode *ConflictResolutionMode         `json:"conflictResolutionMode,omitempty"`
	ImportPrefixes         *[]string                       `json:"importPrefixes,omitempty"`
	MaximumErrors          *int64                          `json:"maximumErrors,omitempty"`
	ProvisioningState      *ImportJobProvisioningStateType `json:"provisioningState,omitempty"`
	Status                 *ImportJobPropertiesStatus      `json:"status,omitempty"`
}
