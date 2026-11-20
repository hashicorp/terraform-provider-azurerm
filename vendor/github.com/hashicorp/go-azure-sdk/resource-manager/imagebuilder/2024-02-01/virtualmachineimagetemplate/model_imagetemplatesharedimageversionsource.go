package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateSource = ImageTemplateSharedImageVersionSource{}

type ImageTemplateSharedImageVersionSource struct {
	ExactVersion   *string `json:"exactVersion,omitempty"`
	ImageVersionId string  `json:"imageVersionId"`

	// Fields inherited from ImageTemplateSource

	Type string `json:"type"`
}

func (s ImageTemplateSharedImageVersionSource) ImageTemplateSource() BaseImageTemplateSourceImpl {
	return BaseImageTemplateSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ImageTemplateSharedImageVersionSource{}

func (s ImageTemplateSharedImageVersionSource) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplateSharedImageVersionSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplateSharedImageVersionSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateSharedImageVersionSource: %+v", err)
	}

	decoded["type"] = "SharedImageVersion"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplateSharedImageVersionSource: %+v", err)
	}

	return encoded, nil
}
