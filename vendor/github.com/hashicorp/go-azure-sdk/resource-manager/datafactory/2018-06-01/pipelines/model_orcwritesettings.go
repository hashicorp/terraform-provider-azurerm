package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FormatWriteSettings = OrcWriteSettings{}

type OrcWriteSettings struct {
	FileNamePrefix *interface{} `json:"fileNamePrefix,omitempty"`
	MaxRowsPerFile *int64       `json:"maxRowsPerFile,omitempty"`

	// Fields inherited from FormatWriteSettings

	Type string `json:"type"`
}

func (s OrcWriteSettings) FormatWriteSettings() BaseFormatWriteSettingsImpl {
	return BaseFormatWriteSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = OrcWriteSettings{}

func (s OrcWriteSettings) MarshalJSON() ([]byte, error) {
	type wrapper OrcWriteSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling OrcWriteSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling OrcWriteSettings: %+v", err)
	}

	decoded["type"] = "OrcWriteSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling OrcWriteSettings: %+v", err)
	}

	return encoded, nil
}
