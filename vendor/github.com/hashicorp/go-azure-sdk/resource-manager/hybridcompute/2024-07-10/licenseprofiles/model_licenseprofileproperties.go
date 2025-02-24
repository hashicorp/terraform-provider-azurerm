package licenseprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LicenseProfileProperties struct {
	EsuProfile        *LicenseProfileArmEsuProperties            `json:"esuProfile,omitempty"`
	ProductProfile    *LicenseProfileArmProductProfileProperties `json:"productProfile,omitempty"`
	ProvisioningState *ProvisioningState                         `json:"provisioningState,omitempty"`
	SoftwareAssurance *LicenseProfilePropertiesSoftwareAssurance `json:"softwareAssurance,omitempty"`
}
