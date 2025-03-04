package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CompressionReadSettings = TarReadSettings{}

type TarReadSettings struct {
	PreserveCompressionFileNameAsFolder *bool `json:"preserveCompressionFileNameAsFolder,omitempty"`

	// Fields inherited from CompressionReadSettings

	Type string `json:"type"`
}

func (s TarReadSettings) CompressionReadSettings() BaseCompressionReadSettingsImpl {
	return BaseCompressionReadSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = TarReadSettings{}

func (s TarReadSettings) MarshalJSON() ([]byte, error) {
	type wrapper TarReadSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TarReadSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TarReadSettings: %+v", err)
	}

	decoded["type"] = "TarReadSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TarReadSettings: %+v", err)
	}

	return encoded, nil
}
