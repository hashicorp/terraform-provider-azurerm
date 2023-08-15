package frontdoors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RouteConfiguration interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawRouteConfigurationImpl struct {
	Type   string
	Values map[string]interface{}
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

	out := RawRouteConfigurationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
