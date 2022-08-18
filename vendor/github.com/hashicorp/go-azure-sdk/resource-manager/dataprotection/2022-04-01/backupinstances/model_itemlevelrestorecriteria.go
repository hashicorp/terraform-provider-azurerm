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

	type RawItemLevelRestoreCriteriaImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawItemLevelRestoreCriteriaImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
