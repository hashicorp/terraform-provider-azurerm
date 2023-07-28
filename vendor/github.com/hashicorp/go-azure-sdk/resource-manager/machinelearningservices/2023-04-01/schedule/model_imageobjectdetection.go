package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutoMLVertical = ImageObjectDetection{}

type ImageObjectDetection struct {
	LimitSettings      *ImageLimitSettings                              `json:"limitSettings,omitempty"`
	ModelSettings      *ImageModelSettingsObjectDetection               `json:"modelSettings,omitempty"`
	PrimaryMetric      *ObjectDetectionPrimaryMetrics                   `json:"primaryMetric,omitempty"`
	SearchSpace        *[]ImageModelDistributionSettingsObjectDetection `json:"searchSpace,omitempty"`
	SweepSettings      *ImageSweepSettings                              `json:"sweepSettings,omitempty"`
	ValidationData     JobInput                                         `json:"validationData"`
	ValidationDataSize *float64                                         `json:"validationDataSize,omitempty"`

	// Fields inherited from AutoMLVertical
	LogVerbosity     *LogVerbosity `json:"logVerbosity,omitempty"`
	TargetColumnName *string       `json:"targetColumnName,omitempty"`
	TrainingData     JobInput      `json:"trainingData"`
}

var _ json.Marshaler = ImageObjectDetection{}

func (s ImageObjectDetection) MarshalJSON() ([]byte, error) {
	type wrapper ImageObjectDetection
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageObjectDetection: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageObjectDetection: %+v", err)
	}
	decoded["taskType"] = "ImageObjectDetection"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageObjectDetection: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &ImageObjectDetection{}

func (s *ImageObjectDetection) UnmarshalJSON(bytes []byte) error {
	type alias ImageObjectDetection
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ImageObjectDetection: %+v", err)
	}

	s.LimitSettings = decoded.LimitSettings
	s.LogVerbosity = decoded.LogVerbosity
	s.ModelSettings = decoded.ModelSettings
	s.PrimaryMetric = decoded.PrimaryMetric
	s.SearchSpace = decoded.SearchSpace
	s.SweepSettings = decoded.SweepSettings
	s.TargetColumnName = decoded.TargetColumnName
	s.ValidationDataSize = decoded.ValidationDataSize

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ImageObjectDetection into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["trainingData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TrainingData' for 'ImageObjectDetection': %+v", err)
		}
		s.TrainingData = impl
	}

	if v, ok := temp["validationData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ValidationData' for 'ImageObjectDetection': %+v", err)
		}
		s.ValidationData = impl
	}
	return nil
}
