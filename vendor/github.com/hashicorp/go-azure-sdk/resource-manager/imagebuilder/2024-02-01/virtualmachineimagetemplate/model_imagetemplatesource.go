package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageTemplateSource interface {
	ImageTemplateSource() BaseImageTemplateSourceImpl
}

var _ ImageTemplateSource = BaseImageTemplateSourceImpl{}

type BaseImageTemplateSourceImpl struct {
	Type string `json:"type"`
}

func (s BaseImageTemplateSourceImpl) ImageTemplateSource() BaseImageTemplateSourceImpl {
	return s
}

var _ ImageTemplateSource = RawImageTemplateSourceImpl{}

// RawImageTemplateSourceImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawImageTemplateSourceImpl struct {
	imageTemplateSource BaseImageTemplateSourceImpl
	Type                string
	Values              map[string]interface{}
}

func (s RawImageTemplateSourceImpl) ImageTemplateSource() BaseImageTemplateSourceImpl {
	return s.imageTemplateSource
}

func UnmarshalImageTemplateSourceImplementation(input []byte) (ImageTemplateSource, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateSource into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ManagedImage") {
		var out ImageTemplateManagedImageSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplateManagedImageSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PlatformImage") {
		var out ImageTemplatePlatformImageSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplatePlatformImageSource: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SharedImageVersion") {
		var out ImageTemplateSharedImageVersionSource
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplateSharedImageVersionSource: %+v", err)
		}
		return out, nil
	}

	var parent BaseImageTemplateSourceImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseImageTemplateSourceImpl: %+v", err)
	}

	return RawImageTemplateSourceImpl{
		imageTemplateSource: parent,
		Type:                value,
		Values:              temp,
	}, nil

}
