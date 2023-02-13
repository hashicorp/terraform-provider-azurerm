package frontdoors

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoutingRuleProperties struct {
	AcceptedProtocols                *[]FrontDoorProtocol                                         `json:"acceptedProtocols,omitempty"`
	EnabledState                     *RoutingRuleEnabledState                                     `json:"enabledState,omitempty"`
	FrontendEndpoints                *[]SubResource                                               `json:"frontendEndpoints,omitempty"`
	PatternsToMatch                  *[]string                                                    `json:"patternsToMatch,omitempty"`
	ResourceState                    *FrontDoorResourceState                                      `json:"resourceState,omitempty"`
	RouteConfiguration               RouteConfiguration                                           `json:"routeConfiguration"`
	RulesEngine                      *SubResource                                                 `json:"rulesEngine,omitempty"`
	WebApplicationFirewallPolicyLink *RoutingRuleUpdateParametersWebApplicationFirewallPolicyLink `json:"webApplicationFirewallPolicyLink,omitempty"`
}

var _ json.Unmarshaler = &RoutingRuleProperties{}

func (s *RoutingRuleProperties) UnmarshalJSON(bytes []byte) error {
	type alias RoutingRuleProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into RoutingRuleProperties: %+v", err)
	}

	s.AcceptedProtocols = decoded.AcceptedProtocols
	s.EnabledState = decoded.EnabledState
	s.FrontendEndpoints = decoded.FrontendEndpoints
	s.PatternsToMatch = decoded.PatternsToMatch
	s.ResourceState = decoded.ResourceState
	s.RulesEngine = decoded.RulesEngine
	s.WebApplicationFirewallPolicyLink = decoded.WebApplicationFirewallPolicyLink

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling RoutingRuleProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["routeConfiguration"]; ok {
		impl, err := unmarshalRouteConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RouteConfiguration' for 'RoutingRuleProperties': %+v", err)
		}
		s.RouteConfiguration = impl
	}
	return nil
}
