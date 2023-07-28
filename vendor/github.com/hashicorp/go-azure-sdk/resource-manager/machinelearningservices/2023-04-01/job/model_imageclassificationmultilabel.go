package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutoMLVertical = ImageClassificationMultilabel{}

type ImageClassificationMultilabel struct {
	LimitSettings      *ImageLimitSettings                             `json:"limitSettings,omitempty"`
	ModelSettings      *ImageModelSettingsClassification               `json:"modelSettings,omitempty"`
	PrimaryMetric      *ClassificationMultilabelPrimaryMetrics         `json:"primaryMetric,omitempty"`
	SearchSpace        *[]ImageModelDistributionSettingsClassification `json:"searchSpace,omitempty"`
	SweepSettings      *ImageSweepSettings                             `json:"sweepSettings,omitempty"`
	ValidationData     JobInput                                        `json:"validationData"`
	ValidationDataSize *float64                                        `json:"validationDataSize,omitempty"`

	// Fields inherited from AutoMLVertical
	LogVerbosity     *LogVerbosity `json:"logVerbosity,omitempty"`
	TargetColumnName *string       `json:"targetColumnName,omitempty"`
	TrainingData     JobInput      `json:"trainingData"`
}

var _ json.Marshaler = ImageClassificationMultilabel{}

func (s ImageClassificationMultilabel) MarshalJSON() ([]byte, error) {
	type wrapper ImageClassificationMultilabel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageClassificationMultilabel: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageClassificationMultilabel: %+v", err)
	}
	decoded["taskType"] = "ImageClassificationMultilabel"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageClassificationMultilabel: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &ImageClassificationMultilabel{}

func (s *ImageClassificationMultilabel) UnmarshalJSON(bytes []byte) error {
	type alias ImageClassificationMultilabel
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ImageClassificationMultilabel: %+v", err)
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
		return fmt.Errorf("unmarshaling ImageClassificationMultilabel into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["trainingData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TrainingData' for 'ImageClassificationMultilabel': %+v", err)
		}
		s.TrainingData = impl
	}

	if v, ok := temp["validationData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ValidationData' for 'ImageClassificationMultilabel': %+v", err)
		}
		s.ValidationData = impl
	}
	return nil
}
