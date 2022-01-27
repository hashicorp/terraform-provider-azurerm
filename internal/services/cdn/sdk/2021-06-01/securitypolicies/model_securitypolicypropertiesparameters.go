package securitypolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

type SecurityPolicyPropertiesParameters interface {
}

func unmarshalSecurityPolicyPropertiesParametersImplementation(input []byte) (SecurityPolicyPropertiesParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SecurityPolicyPropertiesParameters into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "WebApplicationFirewall") {
		var out SecurityPolicyWebApplicationFirewallParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SecurityPolicyWebApplicationFirewallParameters: %+v", err)
		}
		return out, nil
	}

	type RawSecurityPolicyPropertiesParametersImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSecurityPolicyPropertiesParametersImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
