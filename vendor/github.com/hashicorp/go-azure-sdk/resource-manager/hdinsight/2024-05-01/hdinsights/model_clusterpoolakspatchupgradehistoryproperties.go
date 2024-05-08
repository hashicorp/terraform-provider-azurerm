package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterPoolUpgradeHistoryProperties = ClusterPoolAksPatchUpgradeHistoryProperties{}

type ClusterPoolAksPatchUpgradeHistoryProperties struct {
	NewVersion             *string `json:"newVersion,omitempty"`
	OriginalVersion        *string `json:"originalVersion,omitempty"`
	UpgradeAllClusterNodes *bool   `json:"upgradeAllClusterNodes,omitempty"`
	UpgradeClusterPool     *bool   `json:"upgradeClusterPool,omitempty"`

	// Fields inherited from ClusterPoolUpgradeHistoryProperties
	UpgradeResult ClusterPoolUpgradeHistoryUpgradeResultType `json:"upgradeResult"`
	UtcTime       string                                     `json:"utcTime"`
}

var _ json.Marshaler = ClusterPoolAksPatchUpgradeHistoryProperties{}

func (s ClusterPoolAksPatchUpgradeHistoryProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterPoolAksPatchUpgradeHistoryProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterPoolAksPatchUpgradeHistoryProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterPoolAksPatchUpgradeHistoryProperties: %+v", err)
	}
	decoded["upgradeType"] = "AKSPatchUpgrade"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterPoolAksPatchUpgradeHistoryProperties: %+v", err)
	}

	return encoded, nil
}
