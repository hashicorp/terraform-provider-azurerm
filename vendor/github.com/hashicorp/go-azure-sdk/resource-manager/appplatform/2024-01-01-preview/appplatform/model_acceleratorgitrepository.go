package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AcceleratorGitRepository struct {
	AuthSetting       AcceleratorAuthSetting `json:"authSetting"`
	Branch            *string                `json:"branch,omitempty"`
	Commit            *string                `json:"commit,omitempty"`
	GitTag            *string                `json:"gitTag,omitempty"`
	IntervalInSeconds *int64                 `json:"intervalInSeconds,omitempty"`
	SubPath           *string                `json:"subPath,omitempty"`
	Url               string                 `json:"url"`
}

var _ json.Unmarshaler = &AcceleratorGitRepository{}

func (s *AcceleratorGitRepository) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Branch            *string `json:"branch,omitempty"`
		Commit            *string `json:"commit,omitempty"`
		GitTag            *string `json:"gitTag,omitempty"`
		IntervalInSeconds *int64  `json:"intervalInSeconds,omitempty"`
		SubPath           *string `json:"subPath,omitempty"`
		Url               string  `json:"url"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Branch = decoded.Branch
	s.Commit = decoded.Commit
	s.GitTag = decoded.GitTag
	s.IntervalInSeconds = decoded.IntervalInSeconds
	s.SubPath = decoded.SubPath
	s.Url = decoded.Url

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AcceleratorGitRepository into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["authSetting"]; ok {
		impl, err := UnmarshalAcceleratorAuthSettingImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AuthSetting' for 'AcceleratorGitRepository': %+v", err)
		}
		s.AuthSetting = impl
	}

	return nil
}
