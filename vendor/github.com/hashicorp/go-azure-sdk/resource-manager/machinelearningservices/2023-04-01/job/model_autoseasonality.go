package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Seasonality = AutoSeasonality{}

type AutoSeasonality struct {

	// Fields inherited from Seasonality
}

var _ json.Marshaler = AutoSeasonality{}

func (s AutoSeasonality) MarshalJSON() ([]byte, error) {
	type wrapper AutoSeasonality
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutoSeasonality: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutoSeasonality: %+v", err)
	}
	decoded["mode"] = "Auto"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutoSeasonality: %+v", err)
	}

	return encoded, nil
}
