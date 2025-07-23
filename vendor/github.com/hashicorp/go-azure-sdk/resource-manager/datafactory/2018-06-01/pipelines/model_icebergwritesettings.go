package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FormatWriteSettings = IcebergWriteSettings{}

type IcebergWriteSettings struct {

	// Fields inherited from FormatWriteSettings

	Type string `json:"type"`
}

func (s IcebergWriteSettings) FormatWriteSettings() BaseFormatWriteSettingsImpl {
	return BaseFormatWriteSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = IcebergWriteSettings{}

func (s IcebergWriteSettings) MarshalJSON() ([]byte, error) {
	type wrapper IcebergWriteSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling IcebergWriteSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling IcebergWriteSettings: %+v", err)
	}

	decoded["type"] = "IcebergWriteSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling IcebergWriteSettings: %+v", err)
	}

	return encoded, nil
}
