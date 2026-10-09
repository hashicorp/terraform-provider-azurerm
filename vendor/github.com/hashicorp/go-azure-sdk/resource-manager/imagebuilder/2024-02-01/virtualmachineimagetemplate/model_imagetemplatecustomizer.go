package virtualmachineimagetemplate

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageTemplateCustomizer interface {
	ImageTemplateCustomizer() BaseImageTemplateCustomizerImpl
}

var _ ImageTemplateCustomizer = BaseImageTemplateCustomizerImpl{}

type BaseImageTemplateCustomizerImpl struct {
	Name *string `json:"name,omitempty"`
	Type string  `json:"type"`
}

func (s BaseImageTemplateCustomizerImpl) ImageTemplateCustomizer() BaseImageTemplateCustomizerImpl {
	return s
}

var _ ImageTemplateCustomizer = RawImageTemplateCustomizerImpl{}

// RawImageTemplateCustomizerImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawImageTemplateCustomizerImpl struct {
	imageTemplateCustomizer BaseImageTemplateCustomizerImpl
	Type                    string
	Values                  map[string]interface{}
}

func (s RawImageTemplateCustomizerImpl) ImageTemplateCustomizer() BaseImageTemplateCustomizerImpl {
	return s.imageTemplateCustomizer
}

func (s RawImageTemplateCustomizerImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalImageTemplateCustomizerImplementation(input []byte) (ImageTemplateCustomizer, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageTemplateCustomizer into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "File") {
		var out ImageTemplateFileCustomizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplateFileCustomizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PowerShell") {
		var out ImageTemplatePowerShellCustomizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplatePowerShellCustomizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "WindowsRestart") {
		var out ImageTemplateRestartCustomizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplateRestartCustomizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Shell") {
		var out ImageTemplateShellCustomizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplateShellCustomizer: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "WindowsUpdate") {
		var out ImageTemplateWindowsUpdateCustomizer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageTemplateWindowsUpdateCustomizer: %+v", err)
		}
		return out, nil
	}

	var parent BaseImageTemplateCustomizerImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseImageTemplateCustomizerImpl: %+v", err)
	}

	return RawImageTemplateCustomizerImpl{
		imageTemplateCustomizer: parent,
		Type:                    value,
		Values:                  temp,
	}, nil

}
