package securitypolicies

import (
	"encoding/json"
	"fmt"
)

var _ SecurityPolicyPropertiesParameters = SecurityPolicyWebApplicationFirewallParameters{}

type SecurityPolicyWebApplicationFirewallParameters struct {
	Associations *[]SecurityPolicyWebApplicationFirewallAssociation `json:"associations,omitempty"`
	WafPolicy    *ResourceReference                                 `json:"wafPolicy,omitempty"`

	// Fields inherited from SecurityPolicyPropertiesParameters
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
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SecurityPolicyWebApplicationFirewallParameters: %+v", err)
	}
	decoded["type"] = "WebApplicationFirewall"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SecurityPolicyWebApplicationFirewallParameters: %+v", err)
	}

	return encoded, nil
}
