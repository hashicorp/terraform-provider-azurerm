package licenseprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LicenseProfileUpdateProperties struct {
	EsuProfile        *EsuProfileUpdateProperties                      `json:"esuProfile,omitempty"`
	ProductProfile    *ProductProfileUpdateProperties                  `json:"productProfile,omitempty"`
	SoftwareAssurance *LicenseProfileUpdatePropertiesSoftwareAssurance `json:"softwareAssurance,omitempty"`
}
