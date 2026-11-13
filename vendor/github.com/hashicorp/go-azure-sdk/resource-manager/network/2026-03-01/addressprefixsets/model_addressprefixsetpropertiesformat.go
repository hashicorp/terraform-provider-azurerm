package addressprefixsets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddressPrefixSetPropertiesFormat struct {
	AddressPrefixes   []string           `json:"addressPrefixes"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
