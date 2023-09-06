package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetPublicIPAddressConfiguration struct {
	Name       string                                                        `json:"name"`
	Properties *VirtualMachineScaleSetPublicIPAddressConfigurationProperties `json:"properties,omitempty"`
	Sku        *PublicIPAddressSku                                           `json:"sku,omitempty"`
}
