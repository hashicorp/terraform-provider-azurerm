package workspaces

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OutboundRule interface {
	OutboundRule() BaseOutboundRuleImpl
}

var _ OutboundRule = BaseOutboundRuleImpl{}

type BaseOutboundRuleImpl struct {
	Category *RuleCategory `json:"category,omitempty"`
	Status   *RuleStatus   `json:"status,omitempty"`
	Type     RuleType      `json:"type"`
}

func (s BaseOutboundRuleImpl) OutboundRule() BaseOutboundRuleImpl {
	return s
}

var _ OutboundRule = RawOutboundRuleImpl{}

// RawOutboundRuleImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawOutboundRuleImpl struct {
	outboundRule BaseOutboundRuleImpl
	Type         string
	Values       map[string]interface{}
}

func (s RawOutboundRuleImpl) OutboundRule() BaseOutboundRuleImpl {
	return s.outboundRule
}

func UnmarshalOutboundRuleImplementation(input []byte) (OutboundRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling OutboundRule into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "FQDN") {
		var out FqdnOutboundRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FqdnOutboundRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PrivateEndpoint") {
		var out PrivateEndpointOutboundRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PrivateEndpointOutboundRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServiceTag") {
		var out ServiceTagOutboundRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServiceTagOutboundRule: %+v", err)
		}
		return out, nil
	}

	var parent BaseOutboundRuleImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseOutboundRuleImpl: %+v", err)
	}

	return RawOutboundRuleImpl{
		outboundRule: parent,
		Type:         value,
		Values:       temp,
	}, nil

}
