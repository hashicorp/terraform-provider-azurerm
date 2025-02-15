package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FormatWriteSettings = AvroWriteSettings{}

type AvroWriteSettings struct {
	FileNamePrefix  *string `json:"fileNamePrefix,omitempty"`
	MaxRowsPerFile  *int64  `json:"maxRowsPerFile,omitempty"`
	RecordName      *string `json:"recordName,omitempty"`
	RecordNamespace *string `json:"recordNamespace,omitempty"`

	// Fields inherited from FormatWriteSettings

	Type string `json:"type"`
}

func (s AvroWriteSettings) FormatWriteSettings() BaseFormatWriteSettingsImpl {
	return BaseFormatWriteSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = AvroWriteSettings{}

func (s AvroWriteSettings) MarshalJSON() ([]byte, error) {
	type wrapper AvroWriteSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AvroWriteSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AvroWriteSettings: %+v", err)
	}

	decoded["type"] = "AvroWriteSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AvroWriteSettings: %+v", err)
	}

	return encoded, nil
}
