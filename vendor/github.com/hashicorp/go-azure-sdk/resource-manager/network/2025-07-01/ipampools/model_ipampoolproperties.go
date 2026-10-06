package ipampools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPamPoolProperties struct {
	AddressPrefixes   []string           `json:"addressPrefixes"`
	Description       *string            `json:"description,omitempty"`
	DisplayName       *string            `json:"displayName,omitempty"`
	IPAddressType     *[]IPType          `json:"ipAddressType,omitempty"`
	ParentPoolName    *string            `json:"parentPoolName,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
