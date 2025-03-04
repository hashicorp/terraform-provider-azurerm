package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayAuthenticationCertificatePropertiesFormat struct {
	Data              *string            `json:"data,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
