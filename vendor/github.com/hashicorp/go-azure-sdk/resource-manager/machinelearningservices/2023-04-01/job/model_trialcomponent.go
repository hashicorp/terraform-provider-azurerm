package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrialComponent struct {
	CodeId               *string                   `json:"codeId,omitempty"`
	Command              string                    `json:"command"`
	Distribution         DistributionConfiguration `json:"distribution"`
	EnvironmentId        string                    `json:"environmentId"`
	EnvironmentVariables *map[string]string        `json:"environmentVariables,omitempty"`
	Resources            *JobResourceConfiguration `json:"resources,omitempty"`
}

var _ json.Unmarshaler = &TrialComponent{}

func (s *TrialComponent) UnmarshalJSON(bytes []byte) error {
	type alias TrialComponent
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into TrialComponent: %+v", err)
	}

	s.CodeId = decoded.CodeId
	s.Command = decoded.Command
	s.EnvironmentId = decoded.EnvironmentId
	s.EnvironmentVariables = decoded.EnvironmentVariables
	s.Resources = decoded.Resources

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling TrialComponent into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["distribution"]; ok {
		impl, err := unmarshalDistributionConfigurationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Distribution' for 'TrialComponent': %+v", err)
		}
		s.Distribution = impl
	}
	return nil
}
