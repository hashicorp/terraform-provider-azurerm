package firewallpolicyrulecollectiongroups

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FirewallPolicyRule = NatRule{}

type NatRule struct {
	DestinationAddresses *[]string                            `json:"destinationAddresses,omitempty"`
	DestinationPorts     *[]string                            `json:"destinationPorts,omitempty"`
	IPProtocols          *[]FirewallPolicyRuleNetworkProtocol `json:"ipProtocols,omitempty"`
	SourceAddresses      *[]string                            `json:"sourceAddresses,omitempty"`
	SourceIPGroups       *[]string                            `json:"sourceIpGroups,omitempty"`
	TranslatedAddress    *string                              `json:"translatedAddress,omitempty"`
	TranslatedFqdn       *string                              `json:"translatedFqdn,omitempty"`
	TranslatedPort       *string                              `json:"translatedPort,omitempty"`

	// Fields inherited from FirewallPolicyRule
	Description *string `json:"description,omitempty"`
	Name        *string `json:"name,omitempty"`
}

var _ json.Marshaler = NatRule{}

func (s NatRule) MarshalJSON() ([]byte, error) {
	type wrapper NatRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NatRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NatRule: %+v", err)
	}
	decoded["ruleType"] = "NatRule"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NatRule: %+v", err)
	}

	return encoded, nil
}
