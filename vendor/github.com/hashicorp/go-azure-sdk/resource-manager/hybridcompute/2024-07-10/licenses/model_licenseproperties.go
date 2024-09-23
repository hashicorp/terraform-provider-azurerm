package licenses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LicenseProperties struct {
	LicenseDetails    *LicenseDetails    `json:"licenseDetails,omitempty"`
	LicenseType       *LicenseType       `json:"licenseType,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	TenantId          *string            `json:"tenantId,omitempty"`
}
