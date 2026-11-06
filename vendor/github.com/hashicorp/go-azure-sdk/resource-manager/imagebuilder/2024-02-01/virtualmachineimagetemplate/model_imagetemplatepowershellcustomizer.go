package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateCustomizer = ImageTemplatePowerShellCustomizer{}

type ImageTemplatePowerShellCustomizer struct {
	Inline         *[]string `json:"inline,omitempty"`
	RunAsSystem    *bool     `json:"runAsSystem,omitempty"`
	RunElevated    *bool     `json:"runElevated,omitempty"`
	ScriptUri      *string   `json:"scriptUri,omitempty"`
	Sha256Checksum *string   `json:"sha256Checksum,omitempty"`
	ValidExitCodes *[]int64  `json:"validExitCodes,omitempty"`

	// Fields inherited from ImageTemplateCustomizer

	Name *string `json:"name,omitempty"`
	Type string  `json:"type"`
}

func (s ImageTemplatePowerShellCustomizer) ImageTemplateCustomizer() BaseImageTemplateCustomizerImpl {
	return BaseImageTemplateCustomizerImpl{
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = ImageTemplatePowerShellCustomizer{}

func (s ImageTemplatePowerShellCustomizer) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplatePowerShellCustomizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplatePowerShellCustomizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplatePowerShellCustomizer: %+v", err)
	}

	decoded["type"] = "PowerShell"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplatePowerShellCustomizer: %+v", err)
	}

	return encoded, nil
}
