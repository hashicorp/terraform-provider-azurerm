package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FormatWriteSettings = DelimitedTextWriteSettings{}

type DelimitedTextWriteSettings struct {
	FileExtension  string  `json:"fileExtension"`
	FileNamePrefix *string `json:"fileNamePrefix,omitempty"`
	MaxRowsPerFile *int64  `json:"maxRowsPerFile,omitempty"`
	QuoteAllText   *bool   `json:"quoteAllText,omitempty"`

	// Fields inherited from FormatWriteSettings

	Type string `json:"type"`
}

func (s DelimitedTextWriteSettings) FormatWriteSettings() BaseFormatWriteSettingsImpl {
	return BaseFormatWriteSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = DelimitedTextWriteSettings{}

func (s DelimitedTextWriteSettings) MarshalJSON() ([]byte, error) {
	type wrapper DelimitedTextWriteSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DelimitedTextWriteSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DelimitedTextWriteSettings: %+v", err)
	}

	decoded["type"] = "DelimitedTextWriteSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DelimitedTextWriteSettings: %+v", err)
	}

	return encoded, nil
}
