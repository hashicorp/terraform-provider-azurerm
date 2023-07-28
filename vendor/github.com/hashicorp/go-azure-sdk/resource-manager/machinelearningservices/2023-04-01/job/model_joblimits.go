package job

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobLimits interface {
}

func unmarshalJobLimitsImplementation(input []byte) (JobLimits, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling JobLimits into map[string]interface: %+v", err)
	}

	value, ok := temp["jobLimitsType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Command") {
		var out CommandJobLimits
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CommandJobLimits: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Sweep") {
		var out SweepJobLimits
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SweepJobLimits: %+v", err)
		}
		return out, nil
	}

	type RawJobLimitsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawJobLimitsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
