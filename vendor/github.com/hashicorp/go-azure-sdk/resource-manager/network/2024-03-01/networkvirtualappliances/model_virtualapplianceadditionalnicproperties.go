package networkvirtualappliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualApplianceAdditionalNicProperties struct {
	HasPublicIP *bool   `json:"hasPublicIp,omitempty"`
	Name        *string `json:"name,omitempty"`
}
