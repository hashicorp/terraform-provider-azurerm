package endpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeliveryRuleActionParameters interface {
	DeliveryRuleActionParameters() BaseDeliveryRuleActionParametersImpl
}

var _ DeliveryRuleActionParameters = BaseDeliveryRuleActionParametersImpl{}

type BaseDeliveryRuleActionParametersImpl struct {
	TypeName DeliveryRuleActionParametersType `json:"typeName"`
}

func (s BaseDeliveryRuleActionParametersImpl) DeliveryRuleActionParameters() BaseDeliveryRuleActionParametersImpl {
	return s
}

var _ DeliveryRuleActionParameters = RawDeliveryRuleActionParametersImpl{}

// RawDeliveryRuleActionParametersImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDeliveryRuleActionParametersImpl struct {
	deliveryRuleActionParameters BaseDeliveryRuleActionParametersImpl
	Type                         string
	Values                       map[string]interface{}
}

func (s RawDeliveryRuleActionParametersImpl) DeliveryRuleActionParameters() BaseDeliveryRuleActionParametersImpl {
	return s.deliveryRuleActionParameters
}

func UnmarshalDeliveryRuleActionParametersImplementation(input []byte) (DeliveryRuleActionParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleActionParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["typeName"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DeliveryRuleCacheExpirationActionParameters") {
		var out CacheExpirationActionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CacheExpirationActionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleCacheKeyQueryStringBehaviorActionParameters") {
		var out CacheKeyQueryStringActionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CacheKeyQueryStringActionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleHeaderActionParameters") {
		var out HeaderActionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HeaderActionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleOriginGroupOverrideActionParameters") {
		var out OriginGroupOverrideActionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OriginGroupOverrideActionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleRouteConfigurationOverrideActionParameters") {
		var out RouteConfigurationOverrideActionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RouteConfigurationOverrideActionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleUrlRedirectActionParameters") {
		var out URLRedirectActionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into URLRedirectActionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleUrlRewriteActionParameters") {
		var out URLRewriteActionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into URLRewriteActionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleUrlSigningActionParameters") {
		var out URLSigningActionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into URLSigningActionParameters: %+v", err)
		}
		return out, nil
	}

	var parent BaseDeliveryRuleActionParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDeliveryRuleActionParametersImpl: %+v", err)
	}

	return RawDeliveryRuleActionParametersImpl{
		deliveryRuleActionParameters: parent,
		Type:                         value,
		Values:                       temp,
	}, nil

}
