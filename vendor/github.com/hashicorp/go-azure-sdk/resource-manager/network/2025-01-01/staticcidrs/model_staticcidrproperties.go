package staticcidrs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticCidrProperties struct {
	AddressPrefixes               *[]string          `json:"addressPrefixes,omitempty"`
	Description                   *string            `json:"description,omitempty"`
	NumberOfIPAddressesToAllocate *string            `json:"numberOfIPAddressesToAllocate,omitempty"`
	ProvisioningState             *ProvisioningState `json:"provisioningState,omitempty"`
	TotalNumberOfIPAddresses      *string            `json:"totalNumberOfIPAddresses,omitempty"`
}
