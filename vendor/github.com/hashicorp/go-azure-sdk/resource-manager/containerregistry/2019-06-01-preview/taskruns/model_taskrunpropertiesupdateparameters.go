package taskruns

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaskRunPropertiesUpdateParameters struct {
	ForceUpdateTag *string    `json:"forceUpdateTag,omitempty"`
	RunRequest     RunRequest `json:"runRequest"`
}

var _ json.Unmarshaler = &TaskRunPropertiesUpdateParameters{}

func (s *TaskRunPropertiesUpdateParameters) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ForceUpdateTag *string `json:"forceUpdateTag,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ForceUpdateTag = decoded.ForceUpdateTag

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling TaskRunPropertiesUpdateParameters into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["runRequest"]; ok {
		impl, err := UnmarshalRunRequestImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RunRequest' for 'TaskRunPropertiesUpdateParameters': %+v", err)
		}
		s.RunRequest = impl
	}

	return nil
}
