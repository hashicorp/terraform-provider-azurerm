package hdinsights

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ClusterUpgradeHistoryProperties = ClusterAksPatchUpgradeHistoryProperties{}

type ClusterAksPatchUpgradeHistoryProperties struct {
	NewVersion      *string `json:"newVersion,omitempty"`
	OriginalVersion *string `json:"originalVersion,omitempty"`

	// Fields inherited from ClusterUpgradeHistoryProperties
	UpgradeResult ClusterUpgradeHistoryUpgradeResultType `json:"upgradeResult"`
	UtcTime       string                                 `json:"utcTime"`
}

var _ json.Marshaler = ClusterAksPatchUpgradeHistoryProperties{}

func (s ClusterAksPatchUpgradeHistoryProperties) MarshalJSON() ([]byte, error) {
	type wrapper ClusterAksPatchUpgradeHistoryProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClusterAksPatchUpgradeHistoryProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClusterAksPatchUpgradeHistoryProperties: %+v", err)
	}
	decoded["upgradeType"] = "AKSPatchUpgrade"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClusterAksPatchUpgradeHistoryProperties: %+v", err)
	}

	return encoded, nil
}
