package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobService struct {
	Endpoint       *string            `json:"endpoint,omitempty"`
	ErrorMessage   *string            `json:"errorMessage,omitempty"`
	JobServiceType *string            `json:"jobServiceType,omitempty"`
	Nodes          Nodes              `json:"nodes"`
	Port           *int64             `json:"port,omitempty"`
	Properties     *map[string]string `json:"properties,omitempty"`
	Status         *string            `json:"status,omitempty"`
}

var _ json.Unmarshaler = &JobService{}

func (s *JobService) UnmarshalJSON(bytes []byte) error {
	type alias JobService
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into JobService: %+v", err)
	}

	s.Endpoint = decoded.Endpoint
	s.ErrorMessage = decoded.ErrorMessage
	s.JobServiceType = decoded.JobServiceType
	s.Port = decoded.Port
	s.Properties = decoded.Properties
	s.Status = decoded.Status

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling JobService into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["nodes"]; ok {
		impl, err := unmarshalNodesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Nodes' for 'JobService': %+v", err)
		}
		s.Nodes = impl
	}
	return nil
}
