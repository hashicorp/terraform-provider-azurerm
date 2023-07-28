package onlinedeployment

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OnlineScaleSettings = TargetUtilizationScaleSettings{}

type TargetUtilizationScaleSettings struct {
	MaxInstances                *int64  `json:"maxInstances,omitempty"`
	MinInstances                *int64  `json:"minInstances,omitempty"`
	PollingInterval             *string `json:"pollingInterval,omitempty"`
	TargetUtilizationPercentage *int64  `json:"targetUtilizationPercentage,omitempty"`

	// Fields inherited from OnlineScaleSettings
}

var _ json.Marshaler = TargetUtilizationScaleSettings{}

func (s TargetUtilizationScaleSettings) MarshalJSON() ([]byte, error) {
	type wrapper TargetUtilizationScaleSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TargetUtilizationScaleSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TargetUtilizationScaleSettings: %+v", err)
	}
	decoded["scaleType"] = "TargetUtilization"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TargetUtilizationScaleSettings: %+v", err)
	}

	return encoded, nil
}
