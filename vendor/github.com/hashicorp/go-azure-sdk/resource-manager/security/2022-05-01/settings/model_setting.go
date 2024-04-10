package settings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Setting interface {
}

// RawSettingImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSettingImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalSettingImplementation(input []byte) (Setting, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Setting into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AlertSyncSettings") {
		var out AlertSyncSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AlertSyncSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DataExportSettings") {
		var out DataExportSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DataExportSettings: %+v", err)
		}
		return out, nil
	}

	out := RawSettingImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
