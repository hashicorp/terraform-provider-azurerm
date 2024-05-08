package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterUpgradeProperties = ClusterHotfixUpgradeProperties{}

type ClusterHotfixUpgradeProperties struct {
	ComponentName        *string `json:"componentName,omitempty"`
	TargetBuildNumber    *string `json:"targetBuildNumber,omitempty"`
	TargetClusterVersion *string `json:"targetClusterVersion,omitempty"`
	TargetOssVersion     *string `json:"targetOssVersion,omitempty"`

	// Fields inherited from ClusterUpgradeProperties
}

var _ json.Marshaler = ClusterHotfixUpgradeProperties{}

func (s ClusterHotfixUpgradeProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterHotfixUpgradeProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterHotfixUpgradeProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterHotfixUpgradeProperties: %+v", err)
	}
	decoded["upgradeType"] = "HotfixUpgrade"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterHotfixUpgradeProperties: %+v", err)
	}

	return encoded, nil
}
