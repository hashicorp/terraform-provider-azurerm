package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateInVMValidator = ImageTemplateShellValidator{}

type ImageTemplateShellValidator struct {
	Inline         *[]string `json:"inline,omitempty"`
	ScriptUri      *string   `json:"scriptUri,omitempty"`
	Sha256Checksum *string   `json:"sha256Checksum,omitempty"`

	// Fields inherited from ImageTemplateInVMValidator

	Name *string `json:"name,omitempty"`
	Type string  `json:"type"`
}

func (s ImageTemplateShellValidator) ImageTemplateInVMValidator() BaseImageTemplateInVMValidatorImpl {
	return BaseImageTemplateInVMValidatorImpl{
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = ImageTemplateShellValidator{}

func (s ImageTemplateShellValidator) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplateShellValidator
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplateShellValidator: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateShellValidator: %+v", err)
	}

	decoded["type"] = "Shell"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplateShellValidator: %+v", err)
	}

	return encoded, nil
}
