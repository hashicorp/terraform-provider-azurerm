package customdomains

import (
	"encoding/json"
	"fmt"
)

type CustomDomainProperties struct {
	CustomHttpsParameters           CustomDomainHttpsParameters      `json:"customHttpsParameters"`
	CustomHttpsProvisioningState    *CustomHttpsProvisioningState    `json:"customHttpsProvisioningState,omitempty"`
	CustomHttpsProvisioningSubstate *CustomHttpsProvisioningSubstate `json:"customHttpsProvisioningSubstate,omitempty"`
	HostName                        string                           `json:"hostName"`
	ProvisioningState               *string                          `json:"provisioningState,omitempty"`
	ResourceState                   *CustomDomainResourceState       `json:"resourceState,omitempty"`
	ValidationData                  *string                          `json:"validationData,omitempty"`
}

var _ json.Unmarshaler = &CustomDomainProperties{}

func (s *CustomDomainProperties) UnmarshalJSON(bytes []byte) error {
	type alias CustomDomainProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into CustomDomainProperties: %+v", err)
	}

	s.CustomHttpsProvisioningState = decoded.CustomHttpsProvisioningState
	s.CustomHttpsProvisioningSubstate = decoded.CustomHttpsProvisioningSubstate
	s.HostName = decoded.HostName
	s.ProvisioningState = decoded.ProvisioningState
	s.ResourceState = decoded.ResourceState
	s.ValidationData = decoded.ValidationData

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CustomDomainProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["customHttpsParameters"]; ok {
		impl, err := unmarshalCustomDomainHttpsParametersImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CustomHttpsParameters' for 'CustomDomainProperties': %+v", err)
		}
		s.CustomHttpsParameters = impl
	}
	return nil
}
