package taskruns

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaskRunProperties struct {
	ForceUpdateTag    *string            `json:"forceUpdateTag,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	RunRequest        RunRequest         `json:"runRequest"`
	RunResult         *Run               `json:"runResult,omitempty"`
}

var _ json.Unmarshaler = &TaskRunProperties{}

func (s *TaskRunProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ForceUpdateTag    *string            `json:"forceUpdateTag,omitempty"`
		ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
		RunResult         *Run               `json:"runResult,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ForceUpdateTag = decoded.ForceUpdateTag
	s.ProvisioningState = decoded.ProvisioningState
	s.RunResult = decoded.RunResult

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling TaskRunProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["runRequest"]; ok {
		impl, err := UnmarshalRunRequestImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RunRequest' for 'TaskRunProperties': %+v", err)
		}
		s.RunRequest = impl
	}

	return nil
}
