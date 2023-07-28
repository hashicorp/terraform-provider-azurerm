package job

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobInput interface {
}

func unmarshalJobInputImplementation(input []byte) (JobInput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling JobInput into map[string]interface: %+v", err)
	}

	value, ok := temp["jobInputType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "custom_model") {
		var out CustomModelJobInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomModelJobInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "literal") {
		var out LiteralJobInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LiteralJobInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "mlflow_model") {
		var out MLFlowModelJobInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MLFlowModelJobInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "mltable") {
		var out MLTableJobInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MLTableJobInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "triton_model") {
		var out TritonModelJobInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TritonModelJobInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "uri_file") {
		var out UriFileJobInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UriFileJobInput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "uri_folder") {
		var out UriFolderJobInput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UriFolderJobInput: %+v", err)
		}
		return out, nil
	}

	type RawJobInputImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawJobInputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
