package services

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ServicePlacementPolicy interface {
}

func unmarshalServicePlacementPolicyImplementation(input []byte) (ServicePlacementPolicy, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ServicePlacementPolicy into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "InvalidDomain") {
		var out ServicePlacementInvalidDomainPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServicePlacementInvalidDomainPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NonPartiallyPlaceService") {
		var out ServicePlacementNonPartiallyPlaceServicePolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServicePlacementNonPartiallyPlaceServicePolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PreferredPrimaryDomain") {
		var out ServicePlacementPreferPrimaryDomainPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServicePlacementPreferPrimaryDomainPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RequiredDomainDistribution") {
		var out ServicePlacementRequireDomainDistributionPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServicePlacementRequireDomainDistributionPolicy: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RequiredDomain") {
		var out ServicePlacementRequiredDomainPolicy
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServicePlacementRequiredDomainPolicy: %+v", err)
		}
		return out, nil
	}

	type RawServicePlacementPolicyImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawServicePlacementPolicyImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
