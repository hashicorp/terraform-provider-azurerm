package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Seasonality = CustomSeasonality{}

type CustomSeasonality struct {
	Value int64 `json:"value"`

	// Fields inherited from Seasonality
}

var _ json.Marshaler = CustomSeasonality{}

func (s CustomSeasonality) MarshalJSON() ([]byte, error) {
	type wrapper CustomSeasonality
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomSeasonality: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomSeasonality: %+v", err)
	}
	decoded["mode"] = "Custom"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomSeasonality: %+v", err)
	}

	return encoded, nil
}
