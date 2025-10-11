package virtualappliancesites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualApplianceSiteProperties struct {
	AddressPrefix     *string                    `json:"addressPrefix,omitempty"`
	O365Policy        *Office365PolicyProperties `json:"o365Policy,omitempty"`
	ProvisioningState *ProvisioningState         `json:"provisioningState,omitempty"`
}
