package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutoMLVertical = TextNer{}

type TextNer struct {
	FeaturizationSettings *FeaturizationSettings        `json:"featurizationSettings,omitempty"`
	LimitSettings         *NlpVerticalLimitSettings     `json:"limitSettings,omitempty"`
	PrimaryMetric         *ClassificationPrimaryMetrics `json:"primaryMetric,omitempty"`
	ValidationData        JobInput                      `json:"validationData"`

	// Fields inherited from AutoMLVertical
	LogVerbosity     *LogVerbosity `json:"logVerbosity,omitempty"`
	TargetColumnName *string       `json:"targetColumnName,omitempty"`
	TrainingData     JobInput      `json:"trainingData"`
}

var _ json.Marshaler = TextNer{}

func (s TextNer) MarshalJSON() ([]byte, error) {
	type wrapper TextNer
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TextNer: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TextNer: %+v", err)
	}
	decoded["taskType"] = "TextNER"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TextNer: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &TextNer{}

func (s *TextNer) UnmarshalJSON(bytes []byte) error {
	type alias TextNer
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into TextNer: %+v", err)
	}

	s.FeaturizationSettings = decoded.FeaturizationSettings
	s.LimitSettings = decoded.LimitSettings
	s.LogVerbosity = decoded.LogVerbosity
	s.PrimaryMetric = decoded.PrimaryMetric
	s.TargetColumnName = decoded.TargetColumnName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling TextNer into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["trainingData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TrainingData' for 'TextNer': %+v", err)
		}
		s.TrainingData = impl
	}

	if v, ok := temp["validationData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ValidationData' for 'TextNer': %+v", err)
		}
		s.ValidationData = impl
	}
	return nil
}
