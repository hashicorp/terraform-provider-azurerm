package licenseprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LicenseProfileArmEsuProperties struct {
	AssignedLicense            *string         `json:"assignedLicense,omitempty"`
	AssignedLicenseImmutableId *string         `json:"assignedLicenseImmutableId,omitempty"`
	EsuEligibility             *EsuEligibility `json:"esuEligibility,omitempty"`
	EsuKeyState                *EsuKeyState    `json:"esuKeyState,omitempty"`
	EsuKeys                    *[]EsuKey       `json:"esuKeys,omitempty"`
	ServerType                 *EsuServerType  `json:"serverType,omitempty"`
}
