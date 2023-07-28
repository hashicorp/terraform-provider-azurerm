package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutoMLVertical = TextClassificationMultilabel{}

type TextClassificationMultilabel struct {
	FeaturizationSettings *FeaturizationSettings                  `json:"featurizationSettings,omitempty"`
	LimitSettings         *NlpVerticalLimitSettings               `json:"limitSettings,omitempty"`
	PrimaryMetric         *ClassificationMultilabelPrimaryMetrics `json:"primaryMetric,omitempty"`
	ValidationData        JobInput                                `json:"validationData"`

	// Fields inherited from AutoMLVertical
	LogVerbosity     *LogVerbosity `json:"logVerbosity,omitempty"`
	TargetColumnName *string       `json:"targetColumnName,omitempty"`
	TrainingData     JobInput      `json:"trainingData"`
}

var _ json.Marshaler = TextClassificationMultilabel{}

func (s TextClassificationMultilabel) MarshalJSON() ([]byte, error) {
	type wrapper TextClassificationMultilabel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TextClassificationMultilabel: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TextClassificationMultilabel: %+v", err)
	}
	decoded["taskType"] = "TextClassificationMultilabel"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TextClassificationMultilabel: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &TextClassificationMultilabel{}

func (s *TextClassificationMultilabel) UnmarshalJSON(bytes []byte) error {
	type alias TextClassificationMultilabel
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into TextClassificationMultilabel: %+v", err)
	}

	s.FeaturizationSettings = decoded.FeaturizationSettings
	s.LimitSettings = decoded.LimitSettings
	s.LogVerbosity = decoded.LogVerbosity
	s.PrimaryMetric = decoded.PrimaryMetric
	s.TargetColumnName = decoded.TargetColumnName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling TextClassificationMultilabel into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["trainingData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TrainingData' for 'TextClassificationMultilabel': %+v", err)
		}
		s.TrainingData = impl
	}

	if v, ok := temp["validationData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ValidationData' for 'TextClassificationMultilabel': %+v", err)
		}
		s.ValidationData = impl
	}
	return nil
}
