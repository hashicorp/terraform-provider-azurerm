package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LicenseProfileMachineInstanceView struct {
	EsuProfile        *LicenseProfileMachineInstanceViewEsuProperties     `json:"esuProfile,omitempty"`
	LicenseChannel    *string                                             `json:"licenseChannel,omitempty"`
	LicenseStatus     *LicenseStatus                                      `json:"licenseStatus,omitempty"`
	ProductProfile    *LicenseProfileArmProductProfileProperties          `json:"productProfile,omitempty"`
	SoftwareAssurance *LicenseProfileMachineInstanceViewSoftwareAssurance `json:"softwareAssurance,omitempty"`
}
