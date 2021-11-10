package service

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ServiceResourceProperties interface {
}

var _ json.Unmarshaler = &ServiceResourceProperties{}

func (s *ServiceResourceProperties) UnmarshalJSON(bytes []byte) error {
	type alias ServiceResourceProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ServiceResourceProperties: %+v", err)
	}

	s.CorrelationScheme = decoded.CorrelationScheme
	s.DefaultMoveCost = decoded.DefaultMoveCost
	s.PlacementConstraints = decoded.PlacementConstraints
	s.ProvisioningState = decoded.ProvisioningState
	s.ScalingPolicies = decoded.ScalingPolicies
	s.ServiceKind = decoded.ServiceKind
	s.ServiceLoadMetrics = decoded.ServiceLoadMetrics
	s.ServicePackageActivationMode = decoded.ServicePackageActivationMode
	s.ServiceTypeName = decoded.ServiceTypeName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ServiceResourceProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["partitionDescription"]; ok {
		impl, err := unmarshalPartitionImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'PartitionDescription' for 'ServiceResourceProperties': %+v", err)
		}
		s.PartitionDescription = impl
	}

	if v, ok := temp["servicePlacementPolicies"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling ServicePlacementPolicies into list []json.RawMessage: %+v", err)
		}

		output := make([]ServicePlacementPolicy, 0)
		for i, val := range listTemp {
			impl, err := unmarshalServicePlacementPolicyImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'ServicePlacementPolicies' for 'ServiceResourceProperties': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.ServicePlacementPolicies = &output
	}
	return nil
}

func unmarshalServiceResourcePropertiesImplementation(input []byte) (ServiceResourceProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ServiceResourceProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["serviceKind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Stateful") {
		var out StatefulServiceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StatefulServiceProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Stateless") {
		var out StatelessServiceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StatelessServiceProperties: %+v", err)
		}
		return out, nil
	}

	type RawServiceResourcePropertiesImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawServiceResourcePropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
