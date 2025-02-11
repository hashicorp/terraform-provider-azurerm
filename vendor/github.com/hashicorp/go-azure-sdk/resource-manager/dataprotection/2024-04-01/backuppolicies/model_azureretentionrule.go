package backuppolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ BasePolicyRule = AzureRetentionRule{}

type AzureRetentionRule struct {
	IsDefault  *bool             `json:"isDefault,omitempty"`
	Lifecycles []SourceLifeCycle `json:"lifecycles"`

	// Fields inherited from BasePolicyRule

	Name       string `json:"name"`
	ObjectType string `json:"objectType"`
}

func (s AzureRetentionRule) BasePolicyRule() BaseBasePolicyRuleImpl {
	return BaseBasePolicyRuleImpl{
		Name:       s.Name,
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = AzureRetentionRule{}

func (s AzureRetentionRule) MarshalJSON() ([]byte, error) {
	type wrapper AzureRetentionRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureRetentionRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureRetentionRule: %+v", err)
	}

	decoded["objectType"] = "AzureRetentionRule"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureRetentionRule: %+v", err)
	}

	return encoded, nil
}
