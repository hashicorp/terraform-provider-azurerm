package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterPoolAvailableUpgradeProperties = ClusterPoolAvailableUpgradeNodeOsUpgradeProperties{}

type ClusterPoolAvailableUpgradeNodeOsUpgradeProperties struct {
	LatestVersion *string `json:"latestVersion,omitempty"`

	// Fields inherited from ClusterPoolAvailableUpgradeProperties
}

var _ json.Marshaler = ClusterPoolAvailableUpgradeNodeOsUpgradeProperties{}

func (s ClusterPoolAvailableUpgradeNodeOsUpgradeProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterPoolAvailableUpgradeNodeOsUpgradeProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterPoolAvailableUpgradeNodeOsUpgradeProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterPoolAvailableUpgradeNodeOsUpgradeProperties: %+v", err)
	}
	decoded["upgradeType"] = "NodeOsUpgrade"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterPoolAvailableUpgradeNodeOsUpgradeProperties: %+v", err)
	}

	return encoded, nil
}
