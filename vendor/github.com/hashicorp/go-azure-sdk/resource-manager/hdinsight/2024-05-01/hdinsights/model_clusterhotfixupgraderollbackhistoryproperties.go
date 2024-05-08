package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterUpgradeHistoryProperties = ClusterHotfixUpgradeRollbackHistoryProperties{}

type ClusterHotfixUpgradeRollbackHistoryProperties struct {
	ComponentName        *string                            `json:"componentName,omitempty"`
	Severity             *ClusterUpgradeHistorySeverityType `json:"severity,omitempty"`
	SourceBuildNumber    *string                            `json:"sourceBuildNumber,omitempty"`
	SourceClusterVersion *string                            `json:"sourceClusterVersion,omitempty"`
	SourceOssVersion     *string                            `json:"sourceOssVersion,omitempty"`
	TargetBuildNumber    *string                            `json:"targetBuildNumber,omitempty"`
	TargetClusterVersion *string                            `json:"targetClusterVersion,omitempty"`
	TargetOssVersion     *string                            `json:"targetOssVersion,omitempty"`

	// Fields inherited from ClusterUpgradeHistoryProperties
	UpgradeResult ClusterUpgradeHistoryUpgradeResultType `json:"upgradeResult"`
	UtcTime       string                                 `json:"utcTime"`
}

var _ json.Marshaler = ClusterHotfixUpgradeRollbackHistoryProperties{}

func (s ClusterHotfixUpgradeRollbackHistoryProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterHotfixUpgradeRollbackHistoryProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterHotfixUpgradeRollbackHistoryProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterHotfixUpgradeRollbackHistoryProperties: %+v", err)
	}
	decoded["upgradeType"] = "HotfixUpgradeRollback"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterHotfixUpgradeRollbackHistoryProperties: %+v", err)
	}

	return encoded, nil
}
