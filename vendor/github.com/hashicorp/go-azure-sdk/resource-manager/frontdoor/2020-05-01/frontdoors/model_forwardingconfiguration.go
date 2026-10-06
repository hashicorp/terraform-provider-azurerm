package frontdoors

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RouteConfiguration = ForwardingConfiguration{}

type ForwardingConfiguration struct {
	BackendPool          *SubResource                 `json:"backendPool,omitempty"`
	CacheConfiguration   *CacheConfiguration          `json:"cacheConfiguration,omitempty"`
	CustomForwardingPath *string                      `json:"customForwardingPath,omitempty"`
	ForwardingProtocol   *FrontDoorForwardingProtocol `json:"forwardingProtocol,omitempty"`

	// Fields inherited from RouteConfiguration

	OdataType string `json:"@odata.type"`
}

func (s ForwardingConfiguration) RouteConfiguration() BaseRouteConfigurationImpl {
	return BaseRouteConfigurationImpl{
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = ForwardingConfiguration{}

func (s ForwardingConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper ForwardingConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ForwardingConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ForwardingConfiguration: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.FrontDoor.Models.FrontdoorForwardingConfiguration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ForwardingConfiguration: %+v", err)
	}

	return encoded, nil
}
