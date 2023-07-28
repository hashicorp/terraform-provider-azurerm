package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutoMLVertical = ImageInstanceSegmentation{}

type ImageInstanceSegmentation struct {
	LimitSettings      *ImageLimitSettings                              `json:"limitSettings,omitempty"`
	ModelSettings      *ImageModelSettingsObjectDetection               `json:"modelSettings,omitempty"`
	PrimaryMetric      *InstanceSegmentationPrimaryMetrics              `json:"primaryMetric,omitempty"`
	SearchSpace        *[]ImageModelDistributionSettingsObjectDetection `json:"searchSpace,omitempty"`
	SweepSettings      *ImageSweepSettings                              `json:"sweepSettings,omitempty"`
	ValidationData     JobInput                                         `json:"validationData"`
	ValidationDataSize *float64                                         `json:"validationDataSize,omitempty"`

	// Fields inherited from AutoMLVertical
	LogVerbosity     *LogVerbosity `json:"logVerbosity,omitempty"`
	TargetColumnName *string       `json:"targetColumnName,omitempty"`
	TrainingData     JobInput      `json:"trainingData"`
}

var _ json.Marshaler = ImageInstanceSegmentation{}

func (s ImageInstanceSegmentation) MarshalJSON() ([]byte, error) {
	type wrapper ImageInstanceSegmentation
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageInstanceSegmentation: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageInstanceSegmentation: %+v", err)
	}
	decoded["taskType"] = "ImageInstanceSegmentation"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageInstanceSegmentation: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &ImageInstanceSegmentation{}

func (s *ImageInstanceSegmentation) UnmarshalJSON(bytes []byte) error {
	type alias ImageInstanceSegmentation
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ImageInstanceSegmentation: %+v", err)
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
		return fmt.Errorf("unmarshaling ImageInstanceSegmentation into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["trainingData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TrainingData' for 'ImageInstanceSegmentation': %+v", err)
		}
		s.TrainingData = impl
	}

	if v, ok := temp["validationData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ValidationData' for 'ImageInstanceSegmentation': %+v", err)
		}
		s.ValidationData = impl
	}
	return nil
}
