package frontdoors

import (
	"encoding/json"
	"fmt"
	"strings"
)

type RouteConfiguration interface {
}

func unmarshalRouteConfigurationImplementation(input []byte) (RouteConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling RouteConfiguration into map[string]interface: %+v", err)
	}

	value, ok := temp["@odata.type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.FrontDoor.Models.FrontdoorForwardingConfiguration") {
		var out ForwardingConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ForwardingConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.FrontDoor.Models.FrontdoorRedirectConfiguration") {
		var out RedirectConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RedirectConfiguration: %+v", err)
		}
		return out, nil
	}

	type RawRouteConfigurationImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawRouteConfigurationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
