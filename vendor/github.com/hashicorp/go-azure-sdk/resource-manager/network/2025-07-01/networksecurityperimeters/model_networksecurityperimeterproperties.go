package networksecurityperimeters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityPerimeterProperties struct {
	PerimeterGuid     *string               `json:"perimeterGuid,omitempty"`
	ProvisioningState *NspProvisioningState `json:"provisioningState,omitempty"`
}
