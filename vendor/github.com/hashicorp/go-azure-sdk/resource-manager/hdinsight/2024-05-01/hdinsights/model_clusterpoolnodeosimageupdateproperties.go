package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterPoolUpgradeProperties = ClusterPoolNodeOsImageUpdateProperties{}

type ClusterPoolNodeOsImageUpdateProperties struct {

	// Fields inherited from ClusterPoolUpgradeProperties
}

var _ json.Marshaler = ClusterPoolNodeOsImageUpdateProperties{}

func (s ClusterPoolNodeOsImageUpdateProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterPoolNodeOsImageUpdateProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterPoolNodeOsImageUpdateProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterPoolNodeOsImageUpdateProperties: %+v", err)
	}
	decoded["upgradeType"] = "NodeOsUpgrade"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterPoolNodeOsImageUpdateProperties: %+v", err)
	}

	return encoded, nil
}
