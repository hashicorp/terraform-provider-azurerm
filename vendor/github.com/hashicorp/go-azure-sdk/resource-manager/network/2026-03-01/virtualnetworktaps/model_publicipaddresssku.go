package virtualnetworktaps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublicIPAddressSku struct {
	Name *PublicIPAddressSkuName `json:"name,omitempty"`
	Tier *PublicIPAddressSkuTier `json:"tier,omitempty"`
}
