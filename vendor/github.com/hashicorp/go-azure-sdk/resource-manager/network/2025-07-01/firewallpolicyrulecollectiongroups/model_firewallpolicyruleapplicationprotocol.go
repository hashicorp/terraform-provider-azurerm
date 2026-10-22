package firewallpolicyrulecollectiongroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyRuleApplicationProtocol struct {
	Port         *int64                                     `json:"port,omitempty"`
	ProtocolType *FirewallPolicyRuleApplicationProtocolType `json:"protocolType,omitempty"`
}
