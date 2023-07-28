package onlinedeployment

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OnlineScaleSettings = DefaultScaleSettings{}

type DefaultScaleSettings struct {

	// Fields inherited from OnlineScaleSettings
}

var _ json.Marshaler = DefaultScaleSettings{}

func (s DefaultScaleSettings) MarshalJSON() ([]byte, error) {
	type wrapper DefaultScaleSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DefaultScaleSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DefaultScaleSettings: %+v", err)
	}
	decoded["scaleType"] = "Default"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DefaultScaleSettings: %+v", err)
	}

	return encoded, nil
}
