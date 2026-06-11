package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateCustomizer = ImageTemplateFileCustomizer{}

type ImageTemplateFileCustomizer struct {
	Destination    *string `json:"destination,omitempty"`
	Sha256Checksum *string `json:"sha256Checksum,omitempty"`
	SourceUri      *string `json:"sourceUri,omitempty"`

	// Fields inherited from ImageTemplateCustomizer

	Name *string `json:"name,omitempty"`
	Type string  `json:"type"`
}

func (s ImageTemplateFileCustomizer) ImageTemplateCustomizer() BaseImageTemplateCustomizerImpl {
	return BaseImageTemplateCustomizerImpl{
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = ImageTemplateFileCustomizer{}

func (s ImageTemplateFileCustomizer) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplateFileCustomizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplateFileCustomizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateFileCustomizer: %+v", err)
	}

	decoded["type"] = "File"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplateFileCustomizer: %+v", err)
	}

	return encoded, nil
}
