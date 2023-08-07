package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ItemLevelRestoreCriteria interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawItemLevelRestoreCriteriaImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalItemLevelRestoreCriteriaImplementation(input []byte) (ItemLevelRestoreCriteria, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ItemLevelRestoreCriteria into map[string]interface: %+v", err)
	}

	value, ok := temp["objectType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "KubernetesPVRestoreCriteria") {
		var out KubernetesPVRestoreCriteria
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KubernetesPVRestoreCriteria: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KubernetesStorageClassRestoreCriteria") {
		var out KubernetesStorageClassRestoreCriteria
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KubernetesStorageClassRestoreCriteria: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "RangeBasedItemLevelRestoreCriteria") {
		var out RangeBasedItemLevelRestoreCriteria
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into RangeBasedItemLevelRestoreCriteria: %+v", err)
		}
		return out, nil
	}

	out := RawItemLevelRestoreCriteriaImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
