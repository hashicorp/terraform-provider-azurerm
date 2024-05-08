package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterAvailableUpgradeProperties = ClusterAvailableUpgradeHotfixUpgradeProperties{}

type ClusterAvailableUpgradeHotfixUpgradeProperties struct {
	ComponentName        *string   `json:"componentName,omitempty"`
	CreatedTime          *string   `json:"createdTime,omitempty"`
	Description          *string   `json:"description,omitempty"`
	ExtendedProperties   *string   `json:"extendedProperties,omitempty"`
	Severity             *Severity `json:"severity,omitempty"`
	SourceBuildNumber    *string   `json:"sourceBuildNumber,omitempty"`
	SourceClusterVersion *string   `json:"sourceClusterVersion,omitempty"`
	SourceOssVersion     *string   `json:"sourceOssVersion,omitempty"`
	TargetBuildNumber    *string   `json:"targetBuildNumber,omitempty"`
	TargetClusterVersion *string   `json:"targetClusterVersion,omitempty"`
	TargetOssVersion     *string   `json:"targetOssVersion,omitempty"`

	// Fields inherited from ClusterAvailableUpgradeProperties
}

var _ json.Marshaler = ClusterAvailableUpgradeHotfixUpgradeProperties{}

func (s ClusterAvailableUpgradeHotfixUpgradeProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterAvailableUpgradeHotfixUpgradeProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterAvailableUpgradeHotfixUpgradeProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterAvailableUpgradeHotfixUpgradeProperties: %+v", err)
	}
	decoded["upgradeType"] = "HotfixUpgrade"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterAvailableUpgradeHotfixUpgradeProperties: %+v", err)
	}

	return encoded, nil
}
