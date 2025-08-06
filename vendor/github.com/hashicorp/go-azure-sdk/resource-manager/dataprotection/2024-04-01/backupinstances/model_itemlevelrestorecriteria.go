package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ItemLevelRestoreCriteria interface {
	ItemLevelRestoreCriteria() BaseItemLevelRestoreCriteriaImpl
}

var _ ItemLevelRestoreCriteria = BaseItemLevelRestoreCriteriaImpl{}

type BaseItemLevelRestoreCriteriaImpl struct {
	ObjectType string `json:"objectType"`
}

func (s BaseItemLevelRestoreCriteriaImpl) ItemLevelRestoreCriteria() BaseItemLevelRestoreCriteriaImpl {
	return s
}

var _ ItemLevelRestoreCriteria = RawItemLevelRestoreCriteriaImpl{}

// RawItemLevelRestoreCriteriaImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawItemLevelRestoreCriteriaImpl struct {
	itemLevelRestoreCriteria BaseItemLevelRestoreCriteriaImpl
	Type                     string
	Values                   map[string]interface{}
}

func (s RawItemLevelRestoreCriteriaImpl) ItemLevelRestoreCriteria() BaseItemLevelRestoreCriteriaImpl {
	return s.itemLevelRestoreCriteria
}

func UnmarshalItemLevelRestoreCriteriaImplementation(input []byte) (ItemLevelRestoreCriteria, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ItemLevelRestoreCriteria into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ItemPathBasedRestoreCriteria") {
		var out ItemPathBasedRestoreCriteria
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ItemPathBasedRestoreCriteria: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KubernetesClusterRestoreCriteria") {
		var out KubernetesClusterRestoreCriteria
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KubernetesClusterRestoreCriteria: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KubernetesClusterVaultTierRestoreCriteria") {
		var out KubernetesClusterVaultTierRestoreCriteria
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KubernetesClusterVaultTierRestoreCriteria: %+v", err)
		}
		return out, nil
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

	var parent BaseItemLevelRestoreCriteriaImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseItemLevelRestoreCriteriaImpl: %+v", err)
	}

	return RawItemLevelRestoreCriteriaImpl{
		itemLevelRestoreCriteria: parent,
		Type:                     value,
		Values:                   temp,
	}, nil

}
