package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserSourceInfo interface {
	UserSourceInfo() BaseUserSourceInfoImpl
}

var _ UserSourceInfo = BaseUserSourceInfoImpl{}

type BaseUserSourceInfoImpl struct {
	Type    string  `json:"type"`
	Version *string `json:"version,omitempty"`
}

func (s BaseUserSourceInfoImpl) UserSourceInfo() BaseUserSourceInfoImpl {
	return s
}

var _ UserSourceInfo = RawUserSourceInfoImpl{}

// RawUserSourceInfoImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawUserSourceInfoImpl struct {
	userSourceInfo BaseUserSourceInfoImpl
	Type           string
	Values         map[string]interface{}
}

func (s RawUserSourceInfoImpl) UserSourceInfo() BaseUserSourceInfoImpl {
	return s.userSourceInfo
}

func UnmarshalUserSourceInfoImplementation(input []byte) (UserSourceInfo, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling UserSourceInfo into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "BuildResult") {
		var out BuildResultUserSourceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BuildResultUserSourceInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Container") {
		var out CustomContainerUserSourceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomContainerUserSourceInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Jar") {
		var out JarUploadedUserSourceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into JarUploadedUserSourceInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "NetCoreZip") {
		var out NetCoreZipUploadedUserSourceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NetCoreZipUploadedUserSourceInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Source") {
		var out SourceUploadedUserSourceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SourceUploadedUserSourceInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "UploadedUserSourceInfo") {
		var out UploadedUserSourceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UploadedUserSourceInfo: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "War") {
		var out WarUploadedUserSourceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WarUploadedUserSourceInfo: %+v", err)
		}
		return out, nil
	}

	var parent BaseUserSourceInfoImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseUserSourceInfoImpl: %+v", err)
	}

	return RawUserSourceInfoImpl{
		userSourceInfo: parent,
		Type:           value,
		Values:         temp,
	}, nil

}
