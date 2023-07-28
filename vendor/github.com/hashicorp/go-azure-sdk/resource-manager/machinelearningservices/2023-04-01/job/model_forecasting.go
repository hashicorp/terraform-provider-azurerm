package job

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AutoMLVertical = Forecasting{}

type Forecasting struct {
	CvSplitColumnNames    *[]string                           `json:"cvSplitColumnNames,omitempty"`
	FeaturizationSettings *TableVerticalFeaturizationSettings `json:"featurizationSettings,omitempty"`
	ForecastingSettings   *ForecastingSettings                `json:"forecastingSettings,omitempty"`
	LimitSettings         *TableVerticalLimitSettings         `json:"limitSettings,omitempty"`
	NCrossValidations     NCrossValidations                   `json:"nCrossValidations"`
	PrimaryMetric         *ForecastingPrimaryMetrics          `json:"primaryMetric,omitempty"`
	TestData              JobInput                            `json:"testData"`
	TestDataSize          *float64                            `json:"testDataSize,omitempty"`
	TrainingSettings      *ForecastingTrainingSettings        `json:"trainingSettings,omitempty"`
	ValidationData        JobInput                            `json:"validationData"`
	ValidationDataSize    *float64                            `json:"validationDataSize,omitempty"`
	WeightColumnName      *string                             `json:"weightColumnName,omitempty"`

	// Fields inherited from AutoMLVertical
	LogVerbosity     *LogVerbosity `json:"logVerbosity,omitempty"`
	TargetColumnName *string       `json:"targetColumnName,omitempty"`
	TrainingData     JobInput      `json:"trainingData"`
}

var _ json.Marshaler = Forecasting{}

func (s Forecasting) MarshalJSON() ([]byte, error) {
	type wrapper Forecasting
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Forecasting: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Forecasting: %+v", err)
	}
	decoded["taskType"] = "Forecasting"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Forecasting: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &Forecasting{}

func (s *Forecasting) UnmarshalJSON(bytes []byte) error {
	type alias Forecasting
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into Forecasting: %+v", err)
	}

	s.CvSplitColumnNames = decoded.CvSplitColumnNames
	s.FeaturizationSettings = decoded.FeaturizationSettings
	s.ForecastingSettings = decoded.ForecastingSettings
	s.LimitSettings = decoded.LimitSettings
	s.LogVerbosity = decoded.LogVerbosity
	s.PrimaryMetric = decoded.PrimaryMetric
	s.TargetColumnName = decoded.TargetColumnName
	s.TestDataSize = decoded.TestDataSize
	s.TrainingSettings = decoded.TrainingSettings
	s.ValidationDataSize = decoded.ValidationDataSize
	s.WeightColumnName = decoded.WeightColumnName

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling Forecasting into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["nCrossValidations"]; ok {
		impl, err := unmarshalNCrossValidationsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'NCrossValidations' for 'Forecasting': %+v", err)
		}
		s.NCrossValidations = impl
	}

	if v, ok := temp["testData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TestData' for 'Forecasting': %+v", err)
		}
		s.TestData = impl
	}

	if v, ok := temp["trainingData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TrainingData' for 'Forecasting': %+v", err)
		}
		s.TrainingData = impl
	}

	if v, ok := temp["validationData"]; ok {
		impl, err := unmarshalJobInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ValidationData' for 'Forecasting': %+v", err)
		}
		s.ValidationData = impl
	}
	return nil
}
