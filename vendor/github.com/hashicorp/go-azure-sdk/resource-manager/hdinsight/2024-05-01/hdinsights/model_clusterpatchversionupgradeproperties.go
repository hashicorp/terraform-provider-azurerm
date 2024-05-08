package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterUpgradeProperties = ClusterPatchVersionUpgradeProperties{}

type ClusterPatchVersionUpgradeProperties struct {
	ComponentName        *string `json:"componentName,omitempty"`
	TargetBuildNumber    *string `json:"targetBuildNumber,omitempty"`
	TargetClusterVersion *string `json:"targetClusterVersion,omitempty"`
	TargetOssVersion     *string `json:"targetOssVersion,omitempty"`

	// Fields inherited from ClusterUpgradeProperties
}

var _ json.Marshaler = ClusterPatchVersionUpgradeProperties{}

func (s ClusterPatchVersionUpgradeProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterPatchVersionUpgradeProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterPatchVersionUpgradeProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterPatchVersionUpgradeProperties: %+v", err)
	}
	decoded["upgradeType"] = "PatchVersionUpgrade"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterPatchVersionUpgradeProperties: %+v", err)
	}

	return encoded, nil
}
