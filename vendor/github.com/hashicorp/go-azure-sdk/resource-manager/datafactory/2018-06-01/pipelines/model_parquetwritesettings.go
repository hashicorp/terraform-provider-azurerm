package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FormatWriteSettings = ParquetWriteSettings{}

type ParquetWriteSettings struct {
	FileNamePrefix *string `json:"fileNamePrefix,omitempty"`
	MaxRowsPerFile *int64  `json:"maxRowsPerFile,omitempty"`

	// Fields inherited from FormatWriteSettings

	Type string `json:"type"`
}

func (s ParquetWriteSettings) FormatWriteSettings() BaseFormatWriteSettingsImpl {
	return BaseFormatWriteSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ParquetWriteSettings{}

func (s ParquetWriteSettings) MarshalJSON() ([]byte, error) {
	type wrapper ParquetWriteSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ParquetWriteSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ParquetWriteSettings: %+v", err)
	}

	decoded["type"] = "ParquetWriteSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ParquetWriteSettings: %+v", err)
	}

	return encoded, nil
}
