package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FormatWriteSettings = JsonWriteSettings{}

type JsonWriteSettings struct {
	FilePattern *string `json:"filePattern,omitempty"`

	// Fields inherited from FormatWriteSettings

	Type string `json:"type"`
}

func (s JsonWriteSettings) FormatWriteSettings() BaseFormatWriteSettingsImpl {
	return BaseFormatWriteSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = JsonWriteSettings{}

func (s JsonWriteSettings) MarshalJSON() ([]byte, error) {
	type wrapper JsonWriteSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JsonWriteSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JsonWriteSettings: %+v", err)
	}

	decoded["type"] = "JsonWriteSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JsonWriteSettings: %+v", err)
	}

	return encoded, nil
}
