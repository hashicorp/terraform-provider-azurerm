package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageTemplateDistributor interface {
	ImageTemplateDistributor() BaseImageTemplateDistributorImpl
}

var _ ImageTemplateDistributor = BaseImageTemplateDistributorImpl{}

type BaseImageTemplateDistributorImpl struct {
	ArtifactTags  *map[string]string `json:"artifactTags,omitempty"`
	RunOutputName string             `json:"runOutputName"`
	Type          string             `json:"type"`
}

func (s BaseImageTemplateDistributorImpl) ImageTemplateDistributor() BaseImageTemplateDistributorImpl {
	return s
}

var _ ImageTemplateDistributor = RawImageTemplateDistributorImpl{}

// RawImageTemplateDistributorImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawImageTemplateDistributorImpl struct {
	imageTemplateDistributor BaseImageTemplateDistributorImpl
	Type                     string
	Values                   map[string]interface{}
}

func (s RawImageTemplateDistributorImpl) ImageTemplateDistributor() BaseImageTemplateDistributorImpl {
	return s.imageTemplateDistributor
}

func UnmarshalImageTemplateDistributorImplementation(input []byte) (ImageTemplateDistributor, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateDistributor into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ManagedImage") {
		var out ImageTemplateManagedImageDistributor
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplateManagedImageDistributor: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SharedImage") {
		var out ImageTemplateSharedImageDistributor
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplateSharedImageDistributor: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VHD") {
		var out ImageTemplateVhdDistributor
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplateVhdDistributor: %+v", err)
		}
		return out, nil
	}

	var parent BaseImageTemplateDistributorImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseImageTemplateDistributorImpl: %+v", err)
	}

	return RawImageTemplateDistributorImpl{
		imageTemplateDistributor: parent,
		Type:                     value,
		Values:                   temp,
	}, nil

}
