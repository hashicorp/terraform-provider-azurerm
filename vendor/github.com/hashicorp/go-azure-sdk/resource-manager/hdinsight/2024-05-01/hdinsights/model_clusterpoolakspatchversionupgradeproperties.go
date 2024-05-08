package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterPoolUpgradeProperties = ClusterPoolAKSPatchVersionUpgradeProperties{}

type ClusterPoolAKSPatchVersionUpgradeProperties struct {
	TargetAksVersion       *string `json:"targetAksVersion,omitempty"`
	UpgradeAllClusterNodes *bool   `json:"upgradeAllClusterNodes,omitempty"`
	UpgradeClusterPool     *bool   `json:"upgradeClusterPool,omitempty"`

	// Fields inherited from ClusterPoolUpgradeProperties
}

var _ json.Marshaler = ClusterPoolAKSPatchVersionUpgradeProperties{}

func (s ClusterPoolAKSPatchVersionUpgradeProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterPoolAKSPatchVersionUpgradeProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterPoolAKSPatchVersionUpgradeProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterPoolAKSPatchVersionUpgradeProperties: %+v", err)
	}
	decoded["upgradeType"] = "AKSPatchUpgrade"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterPoolAKSPatchVersionUpgradeProperties: %+v", err)
	}

	return encoded, nil
}
