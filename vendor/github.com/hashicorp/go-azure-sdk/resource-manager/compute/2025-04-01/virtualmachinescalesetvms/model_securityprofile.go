package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityProfile struct {
	EncryptionAtHost   *bool               `json:"encryptionAtHost,omitempty"`
	EncryptionIdentity *EncryptionIdentity `json:"encryptionIdentity,omitempty"`
	ProxyAgentSettings *ProxyAgentSettings `json:"proxyAgentSettings,omitempty"`
	SecurityType       *SecurityTypes      `json:"securityType,omitempty"`
	UefiSettings       *UefiSettings       `json:"uefiSettings,omitempty"`
}
