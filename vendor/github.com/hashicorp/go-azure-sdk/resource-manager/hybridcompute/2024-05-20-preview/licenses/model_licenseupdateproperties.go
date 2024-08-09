package licenses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LicenseUpdateProperties struct {
	LicenseDetails *LicenseUpdatePropertiesLicenseDetails `json:"licenseDetails,omitempty"`
	LicenseType    *LicenseType                           `json:"licenseType,omitempty"`
}
