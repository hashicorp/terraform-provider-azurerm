package securitypolicies

import (
	"encoding/json"
	"fmt"
)

type SecurityPolicyUpdateProperties struct {
	Parameters SecurityPolicyPropertiesParameters `json:"parameters"`
}

var _ json.Unmarshaler = &SecurityPolicyUpdateProperties{}

func (s *SecurityPolicyUpdateProperties) UnmarshalJSON(bytes []byte) error {

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SecurityPolicyUpdateProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["parameters"]; ok {
		impl, err := unmarshalSecurityPolicyPropertiesParametersImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Parameters' for 'SecurityPolicyUpdateProperties': %+v", err)
		}
		s.Parameters = impl
	}
	return nil
}
