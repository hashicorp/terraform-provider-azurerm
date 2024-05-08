package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterPoolAvailableUpgradeProperties = ClusterPoolAvailableUpgradeAksPatchUpgradeProperties{}

type ClusterPoolAvailableUpgradeAksPatchUpgradeProperties struct {
	CurrentVersion       *string                             `json:"currentVersion,omitempty"`
	CurrentVersionStatus *CurrentClusterPoolAksVersionStatus `json:"currentVersionStatus,omitempty"`
	LatestVersion        *string                             `json:"latestVersion,omitempty"`

	// Fields inherited from ClusterPoolAvailableUpgradeProperties
}

var _ json.Marshaler = ClusterPoolAvailableUpgradeAksPatchUpgradeProperties{}

func (s ClusterPoolAvailableUpgradeAksPatchUpgradeProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterPoolAvailableUpgradeAksPatchUpgradeProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterPoolAvailableUpgradeAksPatchUpgradeProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterPoolAvailableUpgradeAksPatchUpgradeProperties: %+v", err)
	}
	decoded["upgradeType"] = "AKSPatchUpgrade"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterPoolAvailableUpgradeAksPatchUpgradeProperties: %+v", err)
	}

	return encoded, nil
}
