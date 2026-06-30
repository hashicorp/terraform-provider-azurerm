package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateSource = ImageTemplateManagedImageSource{}

type ImageTemplateManagedImageSource struct {
	ImageId string `json:"imageId"`

	// Fields inherited from ImageTemplateSource

	Type string `json:"type"`
}

func (s ImageTemplateManagedImageSource) ImageTemplateSource() BaseImageTemplateSourceImpl {
	return BaseImageTemplateSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ImageTemplateManagedImageSource{}

func (s ImageTemplateManagedImageSource) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplateManagedImageSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplateManagedImageSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateManagedImageSource: %+v", err)
	}

	decoded["type"] = "ManagedImage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplateManagedImageSource: %+v", err)
	}

	return encoded, nil
}
