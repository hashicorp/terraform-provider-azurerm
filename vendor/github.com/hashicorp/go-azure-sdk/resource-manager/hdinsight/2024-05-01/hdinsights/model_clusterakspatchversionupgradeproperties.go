package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterUpgradeProperties = ClusterAKSPatchVersionUpgradeProperties{}

type ClusterAKSPatchVersionUpgradeProperties struct {

	// Fields inherited from ClusterUpgradeProperties
}

var _ json.Marshaler = ClusterAKSPatchVersionUpgradeProperties{}

func (s ClusterAKSPatchVersionUpgradeProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterAKSPatchVersionUpgradeProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterAKSPatchVersionUpgradeProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterAKSPatchVersionUpgradeProperties: %+v", err)
	}
	decoded["upgradeType"] = "AKSPatchUpgrade"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterAKSPatchVersionUpgradeProperties: %+v", err)
	}

	return encoded, nil
}
