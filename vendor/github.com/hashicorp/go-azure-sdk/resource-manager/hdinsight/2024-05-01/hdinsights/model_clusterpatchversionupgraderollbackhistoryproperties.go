package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterUpgradeHistoryProperties = ClusterPatchVersionUpgradeRollbackHistoryProperties{}

type ClusterPatchVersionUpgradeRollbackHistoryProperties struct {
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

var _ json.Marshaler = ClusterPatchVersionUpgradeRollbackHistoryProperties{}

func (s ClusterPatchVersionUpgradeRollbackHistoryProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterPatchVersionUpgradeRollbackHistoryProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterPatchVersionUpgradeRollbackHistoryProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterPatchVersionUpgradeRollbackHistoryProperties: %+v", err)
	}
	decoded["upgradeType"] = "PatchVersionUpgradeRollback"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterPatchVersionUpgradeRollbackHistoryProperties: %+v", err)
	}

	return encoded, nil
}
