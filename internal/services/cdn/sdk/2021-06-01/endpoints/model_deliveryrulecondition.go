package endpoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DeliveryRuleCondition interface {
}

func unmarshalDeliveryRuleConditionImplementation(input []byte) (DeliveryRuleCondition, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DeliveryRuleCondition into map[string]interface: %+v", err)
	}

	value, ok := temp["name"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "ClientPort") {
		var out DeliveryRuleClientPortCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleClientPortCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Cookies") {
		var out DeliveryRuleCookiesCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleCookiesCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HostName") {
		var out DeliveryRuleHostNameCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleHostNameCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HttpVersion") {
		var out DeliveryRuleHttpVersionCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleHttpVersionCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "IsDevice") {
		var out DeliveryRuleIsDeviceCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleIsDeviceCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PostArgs") {
		var out DeliveryRulePostArgsCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRulePostArgsCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "QueryString") {
		var out DeliveryRuleQueryStringCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleQueryStringCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RemoteAddress") {
		var out DeliveryRuleRemoteAddressCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleRemoteAddressCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RequestBody") {
		var out DeliveryRuleRequestBodyCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleRequestBodyCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RequestHeader") {
		var out DeliveryRuleRequestHeaderCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleRequestHeaderCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RequestMethod") {
		var out DeliveryRuleRequestMethodCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleRequestMethodCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RequestScheme") {
		var out DeliveryRuleRequestSchemeCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleRequestSchemeCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RequestUri") {
		var out DeliveryRuleRequestUriCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleRequestUriCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServerPort") {
		var out DeliveryRuleServerPortCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleServerPortCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SocketAddr") {
		var out DeliveryRuleSocketAddrCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleSocketAddrCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SslProtocol") {
		var out DeliveryRuleSslProtocolCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleSslProtocolCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "UrlFileExtension") {
		var out DeliveryRuleUrlFileExtensionCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleUrlFileExtensionCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "UrlFileName") {
		var out DeliveryRuleUrlFileNameCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleUrlFileNameCondition: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "UrlPath") {
		var out DeliveryRuleUrlPathCondition
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DeliveryRuleUrlPathCondition: %+v", err)
		}
		return out, nil
	}

	type RawDeliveryRuleConditionImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawDeliveryRuleConditionImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
