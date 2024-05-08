package hdinsights

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterUpgradeHistoryProperties interface {
}

// RawClusterUpgradeHistoryPropertiesImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawClusterUpgradeHistoryPropertiesImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalClusterUpgradeHistoryPropertiesImplementation(input []byte) (ClusterUpgradeHistoryProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterUpgradeHistoryProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["upgradeType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AKSPatchUpgrade") {
		var out ClusterAksPatchUpgradeHistoryProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClusterAksPatchUpgradeHistoryProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HotfixUpgrade") {
		var out ClusterHotfixUpgradeHistoryProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClusterHotfixUpgradeHistoryProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HotfixUpgradeRollback") {
		var out ClusterHotfixUpgradeRollbackHistoryProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClusterHotfixUpgradeRollbackHistoryProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PatchVersionUpgrade") {
		var out ClusterPatchVersionUpgradeHistoryProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClusterPatchVersionUpgradeHistoryProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PatchVersionUpgradeRollback") {
		var out ClusterPatchVersionUpgradeRollbackHistoryProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClusterPatchVersionUpgradeRollbackHistoryProperties: %+v", err)
		}
		return out, nil
	}

	out := RawClusterUpgradeHistoryPropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
