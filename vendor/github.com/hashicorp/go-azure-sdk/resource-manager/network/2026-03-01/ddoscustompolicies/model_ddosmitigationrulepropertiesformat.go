package ddoscustompolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DdosMitigationRulePropertiesFormat struct {
	ProvisioningState     *ProvisioningState          `json:"provisioningState,omitempty"`
	SourcePolicyOverrides *[]DdosSourcePolicyOverride `json:"sourcePolicyOverrides,omitempty"`
	TcpDefaultMitigations *DdosTcpDefaultMitigations  `json:"tcpDefaultMitigations,omitempty"`
	TrafficScope          DdosMitigationTrafficScope  `json:"trafficScope"`
	UdpDefaultMitigations *DdosUdpDefaultMitigations  `json:"udpDefaultMitigations,omitempty"`
}
