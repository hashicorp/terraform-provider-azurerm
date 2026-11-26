package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayTrustedClientCertificatePropertiesFormat struct {
	ClientCertIssuerDN *string            `json:"clientCertIssuerDN,omitempty"`
	Data               *string            `json:"data,omitempty"`
	ProvisioningState  *ProvisioningState `json:"provisioningState,omitempty"`
	ValidatedCertData  *string            `json:"validatedCertData,omitempty"`
}
