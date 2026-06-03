package securitypartnerproviders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityPartnerProviderPropertiesFormat struct {
	ConnectionStatus     *SecurityPartnerProviderConnectionStatus `json:"connectionStatus,omitempty"`
	ProvisioningState    *ProvisioningState                       `json:"provisioningState,omitempty"`
	SecurityProviderName *SecurityProviderName                    `json:"securityProviderName,omitempty"`
	VirtualHub           *SubResource                             `json:"virtualHub,omitempty"`
}
