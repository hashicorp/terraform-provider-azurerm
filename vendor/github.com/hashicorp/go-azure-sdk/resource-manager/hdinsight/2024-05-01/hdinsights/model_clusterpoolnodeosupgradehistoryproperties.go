package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterPoolUpgradeHistoryProperties = ClusterPoolNodeOsUpgradeHistoryProperties{}

type ClusterPoolNodeOsUpgradeHistoryProperties struct {
	NewNodeOs *string `json:"newNodeOs,omitempty"`

	// Fields inherited from ClusterPoolUpgradeHistoryProperties
	UpgradeResult ClusterPoolUpgradeHistoryUpgradeResultType `json:"upgradeResult"`
	UtcTime       string                                     `json:"utcTime"`
}

var _ json.Marshaler = ClusterPoolNodeOsUpgradeHistoryProperties{}

func (s ClusterPoolNodeOsUpgradeHistoryProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterPoolNodeOsUpgradeHistoryProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterPoolNodeOsUpgradeHistoryProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterPoolNodeOsUpgradeHistoryProperties: %+v", err)
	}
	decoded["upgradeType"] = "NodeOsUpgrade"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterPoolNodeOsUpgradeHistoryProperties: %+v", err)
	}

	return encoded, nil
}
