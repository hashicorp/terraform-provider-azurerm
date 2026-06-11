package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateSource = ImageTemplatePlatformImageSource{}

type ImageTemplatePlatformImageSource struct {
	ExactVersion *string                    `json:"exactVersion,omitempty"`
	Offer        *string                    `json:"offer,omitempty"`
	PlanInfo     *PlatformImagePurchasePlan `json:"planInfo,omitempty"`
	Publisher    *string                    `json:"publisher,omitempty"`
	Sku          *string                    `json:"sku,omitempty"`
	Version      *string                    `json:"version,omitempty"`

	// Fields inherited from ImageTemplateSource

	Type string `json:"type"`
}

func (s ImageTemplatePlatformImageSource) ImageTemplateSource() BaseImageTemplateSourceImpl {
	return BaseImageTemplateSourceImpl{
		Type: s.Type,
	}
}

var _ json.Marshaler = ImageTemplatePlatformImageSource{}

func (s ImageTemplatePlatformImageSource) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplatePlatformImageSource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplatePlatformImageSource: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplatePlatformImageSource: %+v", err)
	}

	decoded["type"] = "PlatformImage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplatePlatformImageSource: %+v", err)
	}

	return encoded, nil
}
