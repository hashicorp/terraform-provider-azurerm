package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageSweepSettings struct {
	EarlyTermination  EarlyTerminationPolicy `json:"earlyTermination"`
	SamplingAlgorithm SamplingAlgorithmType  `json:"samplingAlgorithm"`
}

var _ json.Unmarshaler = &ImageSweepSettings{}

func (s *ImageSweepSettings) UnmarshalJSON(bytes []byte) error {
	type alias ImageSweepSettings
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ImageSweepSettings: %+v", err)
	}

	s.SamplingAlgorithm = decoded.SamplingAlgorithm

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ImageSweepSettings into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["earlyTermination"]; ok {
		impl, err := unmarshalEarlyTerminationPolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'EarlyTermination' for 'ImageSweepSettings': %+v", err)
		}
		s.EarlyTermination = impl
	}
	return nil
}
