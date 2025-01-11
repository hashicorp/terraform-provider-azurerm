package networkmanagereffectivesecurityadminrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EffectiveBaseSecurityAdminRule = EffectiveSecurityAdminRule{}

type EffectiveSecurityAdminRule struct {
	Properties *AdminPropertiesFormat `json:"properties,omitempty"`

	// Fields inherited from EffectiveBaseSecurityAdminRule

	ConfigurationDescription      *string                            `json:"configurationDescription,omitempty"`
	Id                            *string                            `json:"id,omitempty"`
	Kind                          EffectiveAdminRuleKind             `json:"kind"`
	RuleCollectionAppliesToGroups *[]NetworkManagerSecurityGroupItem `json:"ruleCollectionAppliesToGroups,omitempty"`
	RuleCollectionDescription     *string                            `json:"ruleCollectionDescription,omitempty"`
	RuleGroups                    *[]ConfigurationGroup              `json:"ruleGroups,omitempty"`
}

func (s EffectiveSecurityAdminRule) EffectiveBaseSecurityAdminRule() BaseEffectiveBaseSecurityAdminRuleImpl {
	return BaseEffectiveBaseSecurityAdminRuleImpl{
		ConfigurationDescription:      s.ConfigurationDescription,
		Id:                            s.Id,
		Kind:                          s.Kind,
		RuleCollectionAppliesToGroups: s.RuleCollectionAppliesToGroups,
		RuleCollectionDescription:     s.RuleCollectionDescription,
		RuleGroups:                    s.RuleGroups,
	}
}

var _ json.Marshaler = EffectiveSecurityAdminRule{}

func (s EffectiveSecurityAdminRule) MarshalJSON() ([]byte, error) {
	type wrapper EffectiveSecurityAdminRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EffectiveSecurityAdminRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EffectiveSecurityAdminRule: %+v", err)
	}

	decoded["kind"] = "Custom"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EffectiveSecurityAdminRule: %+v", err)
	}

	return encoded, nil
}
