package frontdoors

import (
	"encoding/json"
	"fmt"
)

type RouteConfiguration interface {
}

func unmarshalRouteConfiguration(body []byte) (RouteConfiguration, error) {
	type intermediateType struct {
		Type string `json:"@odata.type"`
	}
	var intermediate intermediateType
	if err := json.Unmarshal(body, &intermediate); err != nil {
		return nil, fmt.Errorf("decoding: %+v", err)
	}

	switch intermediate.Type {

	case "#Microsoft.Azure.FrontDoor.Models.FrontdoorForwardingConfiguration":
		{
			var out ForwardingConfiguration
			if err := json.Unmarshal(body, &out); err != nil {
				return nil, fmt.Errorf("unmarshalling into %q: %+v", "ForwardingConfiguration", err)
			}
			return &out, nil
		}

	case "#Microsoft.Azure.FrontDoor.Models.FrontdoorRedirectConfiguration":
		{
			var out RedirectConfiguration
			if err := json.Unmarshal(body, &out); err != nil {
				return nil, fmt.Errorf("unmarshalling into %q: %+v", "RedirectConfiguration", err)
			}
			return &out, nil
		}

	}

	return nil, fmt.Errorf("unknown value for OdataType: %q", intermediate.Type)
}
