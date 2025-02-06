package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FileShareConfiguration = SkipFileShareConfiguration{}

type SkipFileShareConfiguration struct {

	// Fields inherited from FileShareConfiguration

	ConfigurationType ConfigurationType `json:"configurationType"`
}

func (s SkipFileShareConfiguration) FileShareConfiguration() BaseFileShareConfigurationImpl {
	return BaseFileShareConfigurationImpl{
		ConfigurationType: s.ConfigurationType,
	}
}

var _ json.Marshaler = SkipFileShareConfiguration{}

func (s SkipFileShareConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper SkipFileShareConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SkipFileShareConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SkipFileShareConfiguration: %+v", err)
	}

	decoded["configurationType"] = "Skip"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SkipFileShareConfiguration: %+v", err)
	}

	return encoded, nil
}
