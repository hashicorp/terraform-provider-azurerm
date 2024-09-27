package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TriggerContext interface {
	TriggerContext() BaseTriggerContextImpl
}

var _ TriggerContext = BaseTriggerContextImpl{}

type BaseTriggerContextImpl struct {
	ObjectType string `json:"objectType"`
}

func (s BaseTriggerContextImpl) TriggerContext() BaseTriggerContextImpl {
	return s
}

var _ TriggerContext = RawTriggerContextImpl{}

// RawTriggerContextImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawTriggerContextImpl struct {
	triggerContext BaseTriggerContextImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawTriggerContextImpl) TriggerContext() BaseTriggerContextImpl {
	return s.triggerContext
}

func UnmarshalTriggerContextImplementation(input []byte) (TriggerContext, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TriggerContext into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AdhocBasedTriggerContext") {
		var out AdhocBasedTriggerContext
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AdhocBasedTriggerContext: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ScheduleBasedTriggerContext") {
		var out ScheduleBasedTriggerContext
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ScheduleBasedTriggerContext: %+v", err)
		}
		return out, nil
	}

	var parent BaseTriggerContextImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseTriggerContextImpl: %+v", err)
	}

	return RawTriggerContextImpl{
		triggerContext: parent,
		Type:           value,
		Values:         temp,
	}, nil

}
