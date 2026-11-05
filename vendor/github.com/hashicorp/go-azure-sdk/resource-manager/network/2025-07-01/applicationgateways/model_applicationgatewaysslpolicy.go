package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewaySslPolicy struct {
	CipherSuites         *[]ApplicationGatewaySslCipherSuite `json:"cipherSuites,omitempty"`
	DisabledSslProtocols *[]ApplicationGatewaySslProtocol    `json:"disabledSslProtocols,omitempty"`
	MinProtocolVersion   *ApplicationGatewaySslProtocol      `json:"minProtocolVersion,omitempty"`
	PolicyName           *ApplicationGatewaySslPolicyName    `json:"policyName,omitempty"`
	PolicyType           *ApplicationGatewaySslPolicyType    `json:"policyType,omitempty"`
}
