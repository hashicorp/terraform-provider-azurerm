package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayBackendSettingsPropertiesFormat struct {
	HostName                       *string                     `json:"hostName,omitempty"`
	PickHostNameFromBackendAddress *bool                       `json:"pickHostNameFromBackendAddress,omitempty"`
	Port                           *int64                      `json:"port,omitempty"`
	Probe                          *SubResource                `json:"probe,omitempty"`
	Protocol                       *ApplicationGatewayProtocol `json:"protocol,omitempty"`
	ProvisioningState              *ProvisioningState          `json:"provisioningState,omitempty"`
	Timeout                        *int64                      `json:"timeout,omitempty"`
	TrustedRootCertificates        *[]SubResource              `json:"trustedRootCertificates,omitempty"`
}
