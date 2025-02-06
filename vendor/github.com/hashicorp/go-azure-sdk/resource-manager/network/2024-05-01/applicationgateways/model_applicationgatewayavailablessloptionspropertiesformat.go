package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayAvailableSslOptionsPropertiesFormat struct {
	AvailableCipherSuites *[]ApplicationGatewaySslCipherSuite `json:"availableCipherSuites,omitempty"`
	AvailableProtocols    *[]ApplicationGatewaySslProtocol    `json:"availableProtocols,omitempty"`
	DefaultPolicy         *ApplicationGatewaySslPolicyName    `json:"defaultPolicy,omitempty"`
	PredefinedPolicies    *[]SubResource                      `json:"predefinedPolicies,omitempty"`
}
