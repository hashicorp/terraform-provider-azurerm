package endpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeliveryRuleAction interface {
	DeliveryRuleAction() BaseDeliveryRuleActionImpl
}

var _ DeliveryRuleAction = BaseDeliveryRuleActionImpl{}

type BaseDeliveryRuleActionImpl struct {
	Name DeliveryRuleActionName `json:"name"`
}

func (s BaseDeliveryRuleActionImpl) DeliveryRuleAction() BaseDeliveryRuleActionImpl {
	return s
}

var _ DeliveryRuleAction = RawDeliveryRuleActionImpl{}

// RawDeliveryRuleActionImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDeliveryRuleActionImpl struct {
	deliveryRuleAction BaseDeliveryRuleActionImpl
	Type               string
	Values             map[string]interface{}
}

func (s RawDeliveryRuleActionImpl) DeliveryRuleAction() BaseDeliveryRuleActionImpl {
	return s.deliveryRuleAction
}

func UnmarshalDeliveryRuleActionImplementation(input []byte) (DeliveryRuleAction, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleAction into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["name"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "CacheExpiration") {
		var out DeliveryRuleCacheExpirationAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleCacheExpirationAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CacheKeyQueryString") {
		var out DeliveryRuleCacheKeyQueryStringAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleCacheKeyQueryStringAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ModifyRequestHeader") {
		var out DeliveryRuleRequestHeaderAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleRequestHeaderAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ModifyResponseHeader") {
		var out DeliveryRuleResponseHeaderAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleResponseHeaderAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RouteConfigurationOverride") {
		var out DeliveryRuleRouteConfigurationOverrideAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleRouteConfigurationOverrideAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OriginGroupOverride") {
		var out OriginGroupOverrideAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OriginGroupOverrideAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "UrlRedirect") {
		var out URLRedirectAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into URLRedirectAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "UrlRewrite") {
		var out URLRewriteAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into URLRewriteAction: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "UrlSigning") {
		var out URLSigningAction
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into URLSigningAction: %+v", err)
		}
		return out, nil
	}

	var parent BaseDeliveryRuleActionImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDeliveryRuleActionImpl: %+v", err)
	}

	return RawDeliveryRuleActionImpl{
		deliveryRuleAction: parent,
		Type:               value,
		Values:             temp,
	}, nil

}
