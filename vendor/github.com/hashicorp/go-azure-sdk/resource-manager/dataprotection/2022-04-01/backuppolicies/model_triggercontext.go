package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TriggerContext interface {
}

func unmarshalTriggerContextImplementation(input []byte) (TriggerContext, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling TriggerContext into map[string]interface: %+v", err)
	}

	value, ok := temp["objectType"].(string)
	if !ok {
		return nil, nil
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

	type RawTriggerContextImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawTriggerContextImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
