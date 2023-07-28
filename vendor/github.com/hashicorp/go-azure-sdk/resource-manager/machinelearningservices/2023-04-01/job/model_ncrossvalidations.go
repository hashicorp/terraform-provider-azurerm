package job

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NCrossValidations interface {
}

func unmarshalNCrossValidationsImplementation(input []byte) (NCrossValidations, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling NCrossValidations into map[string]interface: %+v", err)
	}

	value, ok := temp["mode"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Auto") {
		var out AutoNCrossValidations
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AutoNCrossValidations: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Custom") {
		var out CustomNCrossValidations
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomNCrossValidations: %+v", err)
		}
		return out, nil
	}

	type RawNCrossValidationsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawNCrossValidationsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
