package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LicenseProfileMachineInstanceViewEsuProperties struct {
	AssignedLicense            *License                `json:"assignedLicense,omitempty"`
	AssignedLicenseImmutableId *string                 `json:"assignedLicenseImmutableId,omitempty"`
	EsuEligibility             *EsuEligibility         `json:"esuEligibility,omitempty"`
	EsuKeyState                *EsuKeyState            `json:"esuKeyState,omitempty"`
	EsuKeys                    *[]EsuKey               `json:"esuKeys,omitempty"`
	LicenseAssignmentState     *LicenseAssignmentState `json:"licenseAssignmentState,omitempty"`
	ServerType                 *EsuServerType          `json:"serverType,omitempty"`
}
