package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OSConfiguration = WindowsConfiguration{}

type WindowsConfiguration struct {

	// Fields inherited from OSConfiguration

	OsType OSType `json:"osType"`
}

func (s WindowsConfiguration) OSConfiguration() BaseOSConfigurationImpl {
	return BaseOSConfigurationImpl{
		OsType: s.OsType,
	}
}

var _ json.Marshaler = WindowsConfiguration{}

func (s WindowsConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper WindowsConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WindowsConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WindowsConfiguration: %+v", err)
	}

	decoded["osType"] = "Windows"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WindowsConfiguration: %+v", err)
	}

	return encoded, nil
}
