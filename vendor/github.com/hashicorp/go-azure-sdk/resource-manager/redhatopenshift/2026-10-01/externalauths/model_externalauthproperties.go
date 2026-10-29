package externalauths

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExternalAuthProperties struct {
	Claim             ExternalAuthClaimProfile       `json:"claim"`
	Clients           *[]ExternalAuthClientProfile   `json:"clients,omitempty"`
	Issuer            TokenIssuerProfile             `json:"issuer"`
	ProvisioningState *ExternalAuthProvisioningState `json:"provisioningState,omitempty"`
	Status            *ExternalAuthResourceStatus    `json:"status,omitempty"`
}
