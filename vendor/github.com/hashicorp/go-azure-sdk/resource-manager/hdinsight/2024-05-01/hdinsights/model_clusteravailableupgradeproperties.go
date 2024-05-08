package hdinsights

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterAvailableUpgradeProperties interface {
}

// RawClusterAvailableUpgradePropertiesImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawClusterAvailableUpgradePropertiesImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalClusterAvailableUpgradePropertiesImplementation(input []byte) (ClusterAvailableUpgradeProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterAvailableUpgradeProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["upgradeType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AKSPatchUpgrade") {
		var out ClusterAvailableUpgradeAksPatchUpgradeProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClusterAvailableUpgradeAksPatchUpgradeProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HotfixUpgrade") {
		var out ClusterAvailableUpgradeHotfixUpgradeProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClusterAvailableUpgradeHotfixUpgradeProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PatchVersionUpgrade") {
		var out ClusterAvailableUpgradePatchVersionUpgradeProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClusterAvailableUpgradePatchVersionUpgradeProperties: %+v", err)
		}
		return out, nil
	}

	out := RawClusterAvailableUpgradePropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
