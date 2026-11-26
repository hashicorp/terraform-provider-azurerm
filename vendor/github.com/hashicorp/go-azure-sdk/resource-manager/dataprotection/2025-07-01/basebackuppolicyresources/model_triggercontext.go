package basebackuppolicyresources

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

// RawTriggerContextImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawTriggerContextImpl struct {
	triggerContext BaseTriggerContextImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawTriggerContextImpl) TriggerContext() BaseTriggerContextImpl {
	return s.triggerContext
}

func (s RawTriggerContextImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
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
