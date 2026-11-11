package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleActionParameters = RouteConfigurationOverrideActionParameters{}

type RouteConfigurationOverrideActionParameters struct {
	CacheConfiguration  *CacheConfiguration  `json:"cacheConfiguration,omitempty"`
	OriginGroupOverride *OriginGroupOverride `json:"originGroupOverride,omitempty"`

	// Fields inherited from DeliveryRuleActionParameters

	TypeName DeliveryRuleActionParametersType `json:"typeName"`
}

func (s RouteConfigurationOverrideActionParameters) DeliveryRuleActionParameters() BaseDeliveryRuleActionParametersImpl {
	return BaseDeliveryRuleActionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = RouteConfigurationOverrideActionParameters{}

func (s RouteConfigurationOverrideActionParameters) MarshalJSON() ([]byte, error) {
	type wrapper RouteConfigurationOverrideActionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling RouteConfigurationOverrideActionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling RouteConfigurationOverrideActionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleRouteConfigurationOverrideActionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling RouteConfigurationOverrideActionParameters: %+v", err)
	}

	return encoded, nil
}
