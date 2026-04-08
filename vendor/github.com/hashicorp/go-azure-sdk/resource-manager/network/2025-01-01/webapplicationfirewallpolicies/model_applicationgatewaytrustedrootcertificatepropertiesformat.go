package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayTrustedRootCertificatePropertiesFormat struct {
	Data              *string            `json:"data,omitempty"`
	KeyVaultSecretId  *string            `json:"keyVaultSecretId,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
