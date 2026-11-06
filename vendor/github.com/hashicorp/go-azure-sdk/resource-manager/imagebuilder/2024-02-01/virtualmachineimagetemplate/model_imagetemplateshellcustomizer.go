package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateCustomizer = ImageTemplateShellCustomizer{}

type ImageTemplateShellCustomizer struct {
	Inline         *[]string `json:"inline,omitempty"`
	ScriptUri      *string   `json:"scriptUri,omitempty"`
	Sha256Checksum *string   `json:"sha256Checksum,omitempty"`

	// Fields inherited from ImageTemplateCustomizer

	Name *string `json:"name,omitempty"`
	Type string  `json:"type"`
}

func (s ImageTemplateShellCustomizer) ImageTemplateCustomizer() BaseImageTemplateCustomizerImpl {
	return BaseImageTemplateCustomizerImpl{
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = ImageTemplateShellCustomizer{}

func (s ImageTemplateShellCustomizer) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplateShellCustomizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplateShellCustomizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateShellCustomizer: %+v", err)
	}

	decoded["type"] = "Shell"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplateShellCustomizer: %+v", err)
	}

	return encoded, nil
}
