package networkmanagereffectivesecurityadminrules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ EffectiveBaseSecurityAdminRule = EffectiveDefaultSecurityAdminRule{}

type EffectiveDefaultSecurityAdminRule struct {
	Properties *DefaultAdminPropertiesFormat `json:"properties,omitempty"`

	// Fields inherited from EffectiveBaseSecurityAdminRule
	ConfigurationDescription      *string                            `json:"configurationDescription,omitempty"`
	Id                            *string                            `json:"id,omitempty"`
	RuleCollectionAppliesToGroups *[]NetworkManagerSecurityGroupItem `json:"ruleCollectionAppliesToGroups,omitempty"`
	RuleCollectionDescription     *string                            `json:"ruleCollectionDescription,omitempty"`
	RuleGroups                    *[]ConfigurationGroup              `json:"ruleGroups,omitempty"`
}

var _ json.Marshaler = EffectiveDefaultSecurityAdminRule{}

func (s EffectiveDefaultSecurityAdminRule) MarshalJSON() ([]byte, error) {
	type wrapper EffectiveDefaultSecurityAdminRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EffectiveDefaultSecurityAdminRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EffectiveDefaultSecurityAdminRule: %+v", err)
	}
	decoded["kind"] = "Default"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EffectiveDefaultSecurityAdminRule: %+v", err)
	}

	return encoded, nil
}
