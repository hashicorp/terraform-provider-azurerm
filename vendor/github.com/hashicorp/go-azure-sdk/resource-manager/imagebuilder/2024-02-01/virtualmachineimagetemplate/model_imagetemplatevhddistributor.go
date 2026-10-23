package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateDistributor = ImageTemplateVhdDistributor{}

type ImageTemplateVhdDistributor struct {
	Uri *string `json:"uri,omitempty"`

	// Fields inherited from ImageTemplateDistributor

	ArtifactTags  *map[string]string `json:"artifactTags,omitempty"`
	RunOutputName string             `json:"runOutputName"`
	Type          string             `json:"type"`
}

func (s ImageTemplateVhdDistributor) ImageTemplateDistributor() BaseImageTemplateDistributorImpl {
	return BaseImageTemplateDistributorImpl{
		ArtifactTags:  s.ArtifactTags,
		RunOutputName: s.RunOutputName,
		Type:          s.Type,
	}
}

var _ json.Marshaler = ImageTemplateVhdDistributor{}

func (s ImageTemplateVhdDistributor) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplateVhdDistributor
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplateVhdDistributor: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateVhdDistributor: %+v", err)
	}

	decoded["type"] = "VHD"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplateVhdDistributor: %+v", err)
	}

	return encoded, nil
}
