package endpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeliveryRuleConditionParameters interface {
	DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl
}

var _ DeliveryRuleConditionParameters = BaseDeliveryRuleConditionParametersImpl{}

type BaseDeliveryRuleConditionParametersImpl struct {
	TypeName DeliveryRuleConditionParametersType `json:"typeName"`
}

func (s BaseDeliveryRuleConditionParametersImpl) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return s
}

var _ DeliveryRuleConditionParameters = RawDeliveryRuleConditionParametersImpl{}

// RawDeliveryRuleConditionParametersImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDeliveryRuleConditionParametersImpl struct {
	deliveryRuleConditionParameters BaseDeliveryRuleConditionParametersImpl
	Type                            string
	Values                          map[string]interface{}
}

func (s RawDeliveryRuleConditionParametersImpl) DeliveryRuleConditionParameters() BaseDeliveryRuleConditionParametersImpl {
	return s.deliveryRuleConditionParameters
}

func UnmarshalDeliveryRuleConditionParametersImplementation(input []byte) (DeliveryRuleConditionParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleConditionParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["typeName"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "DeliveryRuleClientPortConditionParameters") {
		var out ClientPortMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClientPortMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleCookiesConditionParameters") {
		var out CookiesMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CookiesMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleHttpVersionConditionParameters") {
		var out HTTPVersionMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HTTPVersionMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleHostNameConditionParameters") {
		var out HostNameMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HostNameMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleIsDeviceConditionParameters") {
		var out IsDeviceMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into IsDeviceMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRulePostArgsConditionParameters") {
		var out PostArgsMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PostArgsMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleQueryStringConditionParameters") {
		var out QueryStringMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into QueryStringMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleRemoteAddressConditionParameters") {
		var out RemoteAddressMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RemoteAddressMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleRequestBodyConditionParameters") {
		var out RequestBodyMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RequestBodyMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleRequestHeaderConditionParameters") {
		var out RequestHeaderMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RequestHeaderMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleRequestMethodConditionParameters") {
		var out RequestMethodMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RequestMethodMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleRequestSchemeConditionParameters") {
		var out RequestSchemeMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RequestSchemeMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleRequestUriConditionParameters") {
		var out RequestUriMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RequestUriMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleServerPortConditionParameters") {
		var out ServerPortMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServerPortMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleSocketAddrConditionParameters") {
		var out SocketAddrMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SocketAddrMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleSslProtocolConditionParameters") {
		var out SslProtocolMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SslProtocolMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleUrlFileExtensionMatchConditionParameters") {
		var out URLFileExtensionMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into URLFileExtensionMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleUrlFilenameConditionParameters") {
		var out URLFileNameMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into URLFileNameMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DeliveryRuleUrlPathMatchConditionParameters") {
		var out URLPathMatchConditionParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into URLPathMatchConditionParameters: %+v", err)
		}
		return out, nil
	}

	var parent BaseDeliveryRuleConditionParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDeliveryRuleConditionParametersImpl: %+v", err)
	}

	return RawDeliveryRuleConditionParametersImpl{
		deliveryRuleConditionParameters: parent,
		Type:                            value,
		Values:                          temp,
	}, nil

}
