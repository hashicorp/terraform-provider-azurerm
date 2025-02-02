package rules

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleAction = DeliveryRuleRouteConfigurationOverrideAction{}

type DeliveryRuleRouteConfigurationOverrideAction struct {
	Parameters RouteConfigurationOverrideActionParameters `json:"parameters"`

	// Fields inherited from DeliveryRuleAction

	Name DeliveryRuleActionName `json:"name"`
}

func (s DeliveryRuleRouteConfigurationOverrideAction) DeliveryRuleAction() BaseDeliveryRuleActionImpl {
	return BaseDeliveryRuleActionImpl{
		Name: s.Name,
	}
}

var _ json.Marshaler = DeliveryRuleRouteConfigurationOverrideAction{}

func (s DeliveryRuleRouteConfigurationOverrideAction) MarshalJSON() ([]byte, error) {
	type wrapper DeliveryRuleRouteConfigurationOverrideAction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DeliveryRuleRouteConfigurationOverrideAction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleRouteConfigurationOverrideAction: %+v", err)
	}

	decoded["name"] = "RouteConfigurationOverride"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DeliveryRuleRouteConfigurationOverrideAction: %+v", err)
	}

	return encoded, nil
}
