package firewallpolicyrulecollectiongroups

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FirewallPolicyRule = ApplicationRule{}

type ApplicationRule struct {
	DestinationAddresses *[]string                                `json:"destinationAddresses,omitempty"`
	FqdnTags             *[]string                                `json:"fqdnTags,omitempty"`
	HTTPHeadersToInsert  *[]FirewallPolicyHTTPHeaderToInsert      `json:"httpHeadersToInsert,omitempty"`
	Protocols            *[]FirewallPolicyRuleApplicationProtocol `json:"protocols,omitempty"`
	SourceAddresses      *[]string                                `json:"sourceAddresses,omitempty"`
	SourceIPGroups       *[]string                                `json:"sourceIpGroups,omitempty"`
	TargetFqdns          *[]string                                `json:"targetFqdns,omitempty"`
	TargetURLs           *[]string                                `json:"targetUrls,omitempty"`
	TerminateTLS         *bool                                    `json:"terminateTLS,omitempty"`
	WebCategories        *[]string                                `json:"webCategories,omitempty"`

	// Fields inherited from FirewallPolicyRule

	Description *string                `json:"description,omitempty"`
	Name        *string                `json:"name,omitempty"`
	RuleType    FirewallPolicyRuleType `json:"ruleType"`
}

func (s ApplicationRule) FirewallPolicyRule() BaseFirewallPolicyRuleImpl {
	return BaseFirewallPolicyRuleImpl{
		Description: s.Description,
		Name:        s.Name,
		RuleType:    s.RuleType,
	}
}

var _ json.Marshaler = ApplicationRule{}

func (s ApplicationRule) MarshalJSON() ([]byte, error) {
	type wrapper ApplicationRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ApplicationRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ApplicationRule: %+v", err)
	}

	decoded["ruleType"] = "ApplicationRule"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ApplicationRule: %+v", err)
	}

	return encoded, nil
}
