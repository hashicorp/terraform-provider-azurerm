package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateInVMValidator = ImageTemplateFileValidator{}

type ImageTemplateFileValidator struct {
	Destination    *string `json:"destination,omitempty"`
	Sha256Checksum *string `json:"sha256Checksum,omitempty"`
	SourceUri      *string `json:"sourceUri,omitempty"`

	// Fields inherited from ImageTemplateInVMValidator

	Name *string `json:"name,omitempty"`
	Type string  `json:"type"`
}

func (s ImageTemplateFileValidator) ImageTemplateInVMValidator() BaseImageTemplateInVMValidatorImpl {
	return BaseImageTemplateInVMValidatorImpl{
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = ImageTemplateFileValidator{}

func (s ImageTemplateFileValidator) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplateFileValidator
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplateFileValidator: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateFileValidator: %+v", err)
	}

	decoded["type"] = "File"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplateFileValidator: %+v", err)
	}

	return encoded, nil
}
