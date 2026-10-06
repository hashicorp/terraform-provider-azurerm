package securitypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SecurityPolicyPropertiesParameters = SecurityPolicyWebApplicationFirewallParameters{}

type SecurityPolicyWebApplicationFirewallParameters struct {
	Associations *[]SecurityPolicyWebApplicationFirewallAssociation `json:"associations,omitempty"`
	WafPolicy    *ResourceReference                                 `json:"wafPolicy,omitempty"`

	// Fields inherited from SecurityPolicyPropertiesParameters

	Type SecurityPolicyType `json:"type"`
}

func (s SecurityPolicyWebApplicationFirewallParameters) SecurityPolicyPropertiesParameters() BaseSecurityPolicyPropertiesParametersImpl {
	return BaseSecurityPolicyPropertiesParametersImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = SecurityPolicyWebApplicationFirewallParameters{}

func (s SecurityPolicyWebApplicationFirewallParameters) MarshalJSON() ([]byte, error) {
	type wrapper SecurityPolicyWebApplicationFirewallParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SecurityPolicyWebApplicationFirewallParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SecurityPolicyWebApplicationFirewallParameters: %+v", err)
	}

	decoded["type"] = "WebApplicationFirewall"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SecurityPolicyWebApplicationFirewallParameters: %+v", err)
	}

	return encoded, nil
}
