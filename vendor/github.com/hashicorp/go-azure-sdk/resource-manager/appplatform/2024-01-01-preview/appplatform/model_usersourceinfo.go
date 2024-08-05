package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserSourceInfo interface {
}

// RawUserSourceInfoImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawUserSourceInfoImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalUserSourceInfoImplementation(input []byte) (UserSourceInfo, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling UserSourceInfo into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
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

	if strings.EqualFold(value, "War") {
		var out WarUploadedUserSourceInfo
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WarUploadedUserSourceInfo: %+v", err)
		}
		return out, nil
	}

	out := RawUserSourceInfoImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
