package forwardingrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForwardingRuleProperties struct {
	DomainName          string               `json:"domainName"`
	ForwardingRuleState *ForwardingRuleState `json:"forwardingRuleState,omitempty"`
	Metadata            *map[string]string   `json:"metadata,omitempty"`
	ProvisioningState   *ProvisioningState   `json:"provisioningState,omitempty"`
	TargetDnsServers    []TargetDnsServer    `json:"targetDnsServers"`
}
