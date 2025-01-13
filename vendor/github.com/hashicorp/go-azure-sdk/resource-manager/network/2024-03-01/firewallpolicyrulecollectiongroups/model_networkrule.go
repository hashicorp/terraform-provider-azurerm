package firewallpolicyrulecollectiongroups

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FirewallPolicyRule = NetworkRule{}

type NetworkRule struct {
	DestinationAddresses *[]string                            `json:"destinationAddresses,omitempty"`
	DestinationFqdns     *[]string                            `json:"destinationFqdns,omitempty"`
	DestinationIPGroups  *[]string                            `json:"destinationIpGroups,omitempty"`
	DestinationPorts     *[]string                            `json:"destinationPorts,omitempty"`
	IPProtocols          *[]FirewallPolicyRuleNetworkProtocol `json:"ipProtocols,omitempty"`
	SourceAddresses      *[]string                            `json:"sourceAddresses,omitempty"`
	SourceIPGroups       *[]string                            `json:"sourceIpGroups,omitempty"`

	// Fields inherited from FirewallPolicyRule

	Description *string                `json:"description,omitempty"`
	Name        *string                `json:"name,omitempty"`
	RuleType    FirewallPolicyRuleType `json:"ruleType"`
}

func (s NetworkRule) FirewallPolicyRule() BaseFirewallPolicyRuleImpl {
	return BaseFirewallPolicyRuleImpl{
		Description: s.Description,
		Name:        s.Name,
		RuleType:    s.RuleType,
	}
}

var _ json.Marshaler = NetworkRule{}

func (s NetworkRule) MarshalJSON() ([]byte, error) {
	type wrapper NetworkRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NetworkRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NetworkRule: %+v", err)
	}

	decoded["ruleType"] = "NetworkRule"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NetworkRule: %+v", err)
	}

	return encoded, nil
}
