package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateCustomizer = ImageTemplateWindowsUpdateCustomizer{}

type ImageTemplateWindowsUpdateCustomizer struct {
	Filters        *[]string `json:"filters,omitempty"`
	SearchCriteria *string   `json:"searchCriteria,omitempty"`
	UpdateLimit    *int64    `json:"updateLimit,omitempty"`

	// Fields inherited from ImageTemplateCustomizer

	Name *string `json:"name,omitempty"`
	Type string  `json:"type"`
}

func (s ImageTemplateWindowsUpdateCustomizer) ImageTemplateCustomizer() BaseImageTemplateCustomizerImpl {
	return BaseImageTemplateCustomizerImpl{
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = ImageTemplateWindowsUpdateCustomizer{}

func (s ImageTemplateWindowsUpdateCustomizer) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplateWindowsUpdateCustomizer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplateWindowsUpdateCustomizer: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateWindowsUpdateCustomizer: %+v", err)
	}

	decoded["type"] = "WindowsUpdate"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplateWindowsUpdateCustomizer: %+v", err)
	}

	return encoded, nil
}
