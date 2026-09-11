package adminrules

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BaseAdminRule interface {
	BaseAdminRule() BaseBaseAdminRuleImpl
}

var _ BaseAdminRule = BaseBaseAdminRuleImpl{}

type BaseBaseAdminRuleImpl struct {
	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Kind       AdminRuleKind          `json:"kind"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

func (s BaseBaseAdminRuleImpl) BaseAdminRule() BaseBaseAdminRuleImpl {
	return s
}

var _ BaseAdminRule = RawBaseAdminRuleImpl{}

// RawBaseAdminRuleImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawBaseAdminRuleImpl struct {
	baseAdminRule BaseBaseAdminRuleImpl
	Type          string
	Values        map[string]interface{}
}

func (s RawBaseAdminRuleImpl) BaseAdminRule() BaseBaseAdminRuleImpl {
	return s.baseAdminRule
}

func (s RawBaseAdminRuleImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalBaseAdminRuleImplementation(input []byte) (BaseAdminRule, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BaseAdminRule into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Custom") {
		var out AdminRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AdminRule: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Default") {
		var out DefaultAdminRule
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DefaultAdminRule: %+v", err)
		}
		return out, nil
	}

	var parent BaseBaseAdminRuleImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseBaseAdminRuleImpl: %+v", err)
	}

	return RawBaseAdminRuleImpl{
		baseAdminRule: parent,
		Type:          value,
		Values:        temp,
	}, nil

}
