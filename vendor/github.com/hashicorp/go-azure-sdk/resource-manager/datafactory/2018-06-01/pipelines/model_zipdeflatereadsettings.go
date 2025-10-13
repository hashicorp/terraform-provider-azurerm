package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CompressionReadSettings = ZipDeflateReadSettings{}

type ZipDeflateReadSettings struct {
	PreserveZipFileNameAsFolder *bool `json:"preserveZipFileNameAsFolder,omitempty"`

	// Fields inherited from CompressionReadSettings

	Type string `json:"type"`
}

func (s ZipDeflateReadSettings) CompressionReadSettings() BaseCompressionReadSettingsImpl {
	return BaseCompressionReadSettingsImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ZipDeflateReadSettings{}

func (s ZipDeflateReadSettings) MarshalJSON() ([]byte, error) {
	type wrapper ZipDeflateReadSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ZipDeflateReadSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ZipDeflateReadSettings: %+v", err)
	}

	decoded["type"] = "ZipDeflateReadSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ZipDeflateReadSettings: %+v", err)
	}

	return encoded, nil
}
