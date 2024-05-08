package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterAvailableUpgradeProperties = ClusterAvailableUpgradeAksPatchUpgradeProperties{}

type ClusterAvailableUpgradeAksPatchUpgradeProperties struct {
	CurrentVersion       *string                         `json:"currentVersion,omitempty"`
	CurrentVersionStatus *CurrentClusterAksVersionStatus `json:"currentVersionStatus,omitempty"`
	LatestVersion        *string                         `json:"latestVersion,omitempty"`

	// Fields inherited from ClusterAvailableUpgradeProperties
}

var _ json.Marshaler = ClusterAvailableUpgradeAksPatchUpgradeProperties{}

func (s ClusterAvailableUpgradeAksPatchUpgradeProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterAvailableUpgradeAksPatchUpgradeProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterAvailableUpgradeAksPatchUpgradeProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterAvailableUpgradeAksPatchUpgradeProperties: %+v", err)
	}
	decoded["upgradeType"] = "AKSPatchUpgrade"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterAvailableUpgradeAksPatchUpgradeProperties: %+v", err)
	}

	return encoded, nil
}
