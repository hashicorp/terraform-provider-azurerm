package securitypolicies

import (
	"encoding/json"
	"fmt"
)

type SecurityPolicyProperties struct {
	DeploymentStatus  *DeploymentStatus                  `json:"deploymentStatus,omitempty"`
	Parameters        SecurityPolicyPropertiesParameters `json:"parameters"`
	ProfileName       *string                            `json:"profileName,omitempty"`
	ProvisioningState *AfdProvisioningState              `json:"provisioningState,omitempty"`
}

var _ json.Unmarshaler = &SecurityPolicyProperties{}

func (s *SecurityPolicyProperties) UnmarshalJSON(bytes []byte) error {
	type alias SecurityPolicyProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into SecurityPolicyProperties: %+v", err)
	}

	s.DeploymentStatus = decoded.DeploymentStatus
	s.ProfileName = decoded.ProfileName
	s.ProvisioningState = decoded.ProvisioningState

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SecurityPolicyProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["parameters"]; ok {
		impl, err := unmarshalSecurityPolicyPropertiesParametersImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Parameters' for 'SecurityPolicyProperties': %+v", err)
		}
		s.Parameters = impl
	}
	return nil
}
