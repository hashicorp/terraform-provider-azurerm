package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ImageTemplateDistributor = ImageTemplateManagedImageDistributor{}

type ImageTemplateManagedImageDistributor struct {
	ImageId  string `json:"imageId"`
	Location string `json:"location"`

	// Fields inherited from ImageTemplateDistributor

	ArtifactTags  *map[string]string `json:"artifactTags,omitempty"`
	RunOutputName string             `json:"runOutputName"`
	Type          string             `json:"type"`
}

func (s ImageTemplateManagedImageDistributor) ImageTemplateDistributor() BaseImageTemplateDistributorImpl {
	return BaseImageTemplateDistributorImpl{
		ArtifactTags:  s.ArtifactTags,
		RunOutputName: s.RunOutputName,
		Type:          s.Type,
	}
}

var _ json.Marshaler = ImageTemplateManagedImageDistributor{}

func (s ImageTemplateManagedImageDistributor) MarshalJSON() ([]byte, error) {
	type wrapper ImageTemplateManagedImageDistributor
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageTemplateManagedImageDistributor: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateManagedImageDistributor: %+v", err)
	}

	decoded["type"] = "ManagedImage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageTemplateManagedImageDistributor: %+v", err)
	}

	return encoded, nil
}
