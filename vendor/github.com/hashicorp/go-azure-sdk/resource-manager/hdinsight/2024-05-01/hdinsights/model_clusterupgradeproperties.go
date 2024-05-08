package hdinsights

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterUpgradeProperties interface {
}

// RawClusterUpgradePropertiesImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawClusterUpgradePropertiesImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalClusterUpgradePropertiesImplementation(input []byte) (ClusterUpgradeProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterUpgradeProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["upgradeType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AKSPatchUpgrade") {
		var out ClusterAKSPatchVersionUpgradeProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClusterAKSPatchVersionUpgradeProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HotfixUpgrade") {
		var out ClusterHotfixUpgradeProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClusterHotfixUpgradeProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PatchVersionUpgrade") {
		var out ClusterPatchVersionUpgradeProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClusterPatchVersionUpgradeProperties: %+v", err)
		}
		return out, nil
	}

	out := RawClusterUpgradePropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
