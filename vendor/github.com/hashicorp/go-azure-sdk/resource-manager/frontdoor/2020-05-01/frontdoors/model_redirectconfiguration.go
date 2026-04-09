package frontdoors

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ RouteConfiguration = RedirectConfiguration{}

type RedirectConfiguration struct {
	CustomFragment    *string                    `json:"customFragment,omitempty"`
	CustomHost        *string                    `json:"customHost,omitempty"`
	CustomPath        *string                    `json:"customPath,omitempty"`
	CustomQueryString *string                    `json:"customQueryString,omitempty"`
	RedirectProtocol  *FrontDoorRedirectProtocol `json:"redirectProtocol,omitempty"`
	RedirectType      *FrontDoorRedirectType     `json:"redirectType,omitempty"`

	// Fields inherited from RouteConfiguration

	OdataType string `json:"@odata.type"`
}

func (s RedirectConfiguration) RouteConfiguration() BaseRouteConfigurationImpl {
	return BaseRouteConfigurationImpl{
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = RedirectConfiguration{}

func (s RedirectConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper RedirectConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RedirectConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RedirectConfiguration: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.FrontDoor.Models.FrontdoorRedirectConfiguration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RedirectConfiguration: %+v", err)
	}

	return encoded, nil
}
