package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewaySslCertificatePropertiesFormat struct {
	Data              *string                       `json:"data,omitempty"`
	Hsm               *ApplicationGatewayManagedHsm `json:"hsm,omitempty"`
	KeyVaultSecretId  *string                       `json:"keyVaultSecretId,omitempty"`
	Password          *string                       `json:"password,omitempty"`
	ProvisioningState *ProvisioningState            `json:"provisioningState,omitempty"`
	PublicCertData    *string                       `json:"publicCertData,omitempty"`
}
