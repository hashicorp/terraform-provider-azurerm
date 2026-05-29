package restorepoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorePointSourceMetadata struct {
	DiagnosticsProfile *DiagnosticsProfile                 `json:"diagnosticsProfile,omitempty"`
	HardwareProfile    *HardwareProfile                    `json:"hardwareProfile,omitempty"`
	HyperVGeneration   *HyperVGenerationTypes              `json:"hyperVGeneration,omitempty"`
	LicenseType        *string                             `json:"licenseType,omitempty"`
	Location           *string                             `json:"location,omitempty"`
	OsProfile          *OSProfile                          `json:"osProfile,omitempty"`
	SecurityProfile    *SecurityProfile                    `json:"securityProfile,omitempty"`
	StorageProfile     *RestorePointSourceVMStorageProfile `json:"storageProfile,omitempty"`
	UserData           *string                             `json:"userData,omitempty"`
	VMId               *string                             `json:"vmId,omitempty"`
}
