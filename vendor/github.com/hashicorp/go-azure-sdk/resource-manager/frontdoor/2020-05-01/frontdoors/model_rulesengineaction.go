package frontdoors

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RulesEngineAction struct {
	RequestHeaderActions       *[]HeaderAction    `json:"requestHeaderActions,omitempty"`
	ResponseHeaderActions      *[]HeaderAction    `json:"responseHeaderActions,omitempty"`
	RouteConfigurationOverride RouteConfiguration `json:"routeConfigurationOverride"`
}

var _ json.Unmarshaler = &RulesEngineAction{}

func (s *RulesEngineAction) UnmarshalJSON(bytes []byte) error {
	type alias RulesEngineAction
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into RulesEngineAction: %+v", err)
	}

	s.RequestHeaderActions = decoded.RequestHeaderActions
	s.ResponseHeaderActions = decoded.ResponseHeaderActions

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling RulesEngineAction into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["routeConfigurationOverride"]; ok {
		impl, err := unmarshalRouteConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RouteConfigurationOverride' for 'RulesEngineAction': %+v", err)
		}
		s.RouteConfigurationOverride = impl
	}
	return nil
}
