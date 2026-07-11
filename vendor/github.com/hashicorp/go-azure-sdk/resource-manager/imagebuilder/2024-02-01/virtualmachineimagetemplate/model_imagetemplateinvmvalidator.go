package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageTemplateInVMValidator interface {
	ImageTemplateInVMValidator() BaseImageTemplateInVMValidatorImpl
}

var _ ImageTemplateInVMValidator = BaseImageTemplateInVMValidatorImpl{}

type BaseImageTemplateInVMValidatorImpl struct {
	Name *string `json:"name,omitempty"`
	Type string  `json:"type"`
}

func (s BaseImageTemplateInVMValidatorImpl) ImageTemplateInVMValidator() BaseImageTemplateInVMValidatorImpl {
	return s
}

var _ ImageTemplateInVMValidator = RawImageTemplateInVMValidatorImpl{}

// RawImageTemplateInVMValidatorImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawImageTemplateInVMValidatorImpl struct {
	imageTemplateInVMValidator BaseImageTemplateInVMValidatorImpl
	Type                       string
	Values                     map[string]interface{}
}

func (s RawImageTemplateInVMValidatorImpl) ImageTemplateInVMValidator() BaseImageTemplateInVMValidatorImpl {
	return s.imageTemplateInVMValidator
}

func UnmarshalImageTemplateInVMValidatorImplementation(input []byte) (ImageTemplateInVMValidator, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateInVMValidator into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "File") {
		var out ImageTemplateFileValidator
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplateFileValidator: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PowerShell") {
		var out ImageTemplatePowerShellValidator
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplatePowerShellValidator: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Shell") {
		var out ImageTemplateShellValidator
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplateShellValidator: %+v", err)
		}
		return out, nil
	}

	var parent BaseImageTemplateInVMValidatorImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseImageTemplateInVMValidatorImpl: %+v", err)
	}

	return RawImageTemplateInVMValidatorImpl{
		imageTemplateInVMValidator: parent,
		Type:                       value,
		Values:                     temp,
	}, nil

}
