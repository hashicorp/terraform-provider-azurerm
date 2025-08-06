package networksecurityperimeterconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityPerimeterConfigurationProperties struct {
	NetworkSecurityPerimeter *NetworkSecurityPerimeter        `json:"networkSecurityPerimeter,omitempty"`
	Profile                  *NetworkSecurityPerimeterProfile `json:"profile,omitempty"`
	ProvisioningIssues       *[]ProvisioningIssue             `json:"provisioningIssues,omitempty"`
	ProvisioningState        *string                          `json:"provisioningState,omitempty"`
	ResourceAssociation      *ResourceAssociation             `json:"resourceAssociation,omitempty"`
}
