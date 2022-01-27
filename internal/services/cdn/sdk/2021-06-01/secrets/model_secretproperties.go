package secrets

import (
	"encoding/json"
	"fmt"
)

type SecretProperties struct {
	DeploymentStatus  *DeploymentStatus     `json:"deploymentStatus,omitempty"`
	Parameters        SecretParameters      `json:"parameters"`
	ProfileName       *string               `json:"profileName,omitempty"`
	ProvisioningState *AfdProvisioningState `json:"provisioningState,omitempty"`
}

var _ json.Unmarshaler = &SecretProperties{}

func (s *SecretProperties) UnmarshalJSON(bytes []byte) error {
	type alias SecretProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into SecretProperties: %+v", err)
	}

	s.DeploymentStatus = decoded.DeploymentStatus
	s.ProfileName = decoded.ProfileName
	s.ProvisioningState = decoded.ProvisioningState

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SecretProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["parameters"]; ok {
		impl, err := unmarshalSecretParametersImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Parameters' for 'SecretProperties': %+v", err)
		}
		s.Parameters = impl
	}
	return nil
}
