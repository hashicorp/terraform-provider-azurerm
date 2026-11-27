package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateCustomizer = ImageTemplateRestartCustomizer{}

type ImageTemplateRestartCustomizer struct {
	RestartCheckCommand *string `json:"restartCheckCommand,omitempty"`
	RestartCommand      *string `json:"restartCommand,omitempty"`
	RestartTimeout      *string `json:"restartTimeout,omitempty"`

	// Fields inherited from ImageTemplateCustomizer

	Name *string `json:"name,omitempty"`
	Type string  `json:"type"`
}

func (s ImageTemplateRestartCustomizer) ImageTemplateCustomizer() BaseImageTemplateCustomizerImpl {
	return BaseImageTemplateCustomizerImpl{
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = ImageTemplateRestartCustomizer{}

func (s ImageTemplateRestartCustomizer) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplateRestartCustomizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplateRestartCustomizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateRestartCustomizer: %+v", err)
	}

	decoded["type"] = "WindowsRestart"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplateRestartCustomizer: %+v", err)
	}

	return encoded, nil
}
