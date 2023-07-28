package job

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobOutput interface {
}

func unmarshalJobOutputImplementation(input []byte) (JobOutput, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling JobOutput into map[string]interface: %+v", err)
	}

	value, ok := temp["jobOutputType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "custom_model") {
		var out CustomModelJobOutput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomModelJobOutput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "mlflow_model") {
		var out MLFlowModelJobOutput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MLFlowModelJobOutput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "mltable") {
		var out MLTableJobOutput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MLTableJobOutput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "triton_model") {
		var out TritonModelJobOutput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TritonModelJobOutput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "uri_file") {
		var out UriFileJobOutput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UriFileJobOutput: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "uri_folder") {
		var out UriFolderJobOutput
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into UriFolderJobOutput: %+v", err)
		}
		return out, nil
	}

	type RawJobOutputImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawJobOutputImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
