package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ForecastingSettings struct {
	CountryOrRegionForHolidays *string                           `json:"countryOrRegionForHolidays,omitempty"`
	CvStepSize                 *int64                            `json:"cvStepSize,omitempty"`
	FeatureLags                *FeatureLags                      `json:"featureLags,omitempty"`
	ForecastHorizon            ForecastHorizon                   `json:"forecastHorizon"`
	Frequency                  *string                           `json:"frequency,omitempty"`
	Seasonality                Seasonality                       `json:"seasonality"`
	ShortSeriesHandlingConfig  *ShortSeriesHandlingConfiguration `json:"shortSeriesHandlingConfig,omitempty"`
	TargetAggregateFunction    *TargetAggregationFunction        `json:"targetAggregateFunction,omitempty"`
	TargetLags                 TargetLags                        `json:"targetLags"`
	TargetRollingWindowSize    TargetRollingWindowSize           `json:"targetRollingWindowSize"`
	TimeColumnName             *string                           `json:"timeColumnName,omitempty"`
	TimeSeriesIdColumnNames    *[]string                         `json:"timeSeriesIdColumnNames,omitempty"`
	UseStl                     *UseStl                           `json:"useStl,omitempty"`
}

var _ json.Unmarshaler = &ForecastingSettings{}

func (s *ForecastingSettings) UnmarshalJSON(bytes []byte) error {
	type alias ForecastingSettings
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ForecastingSettings: %+v", err)
	}

	s.CountryOrRegionForHolidays = decoded.CountryOrRegionForHolidays
	s.CvStepSize = decoded.CvStepSize
	s.FeatureLags = decoded.FeatureLags
	s.Frequency = decoded.Frequency
	s.ShortSeriesHandlingConfig = decoded.ShortSeriesHandlingConfig
	s.TargetAggregateFunction = decoded.TargetAggregateFunction
	s.TimeColumnName = decoded.TimeColumnName
	s.TimeSeriesIdColumnNames = decoded.TimeSeriesIdColumnNames
	s.UseStl = decoded.UseStl

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ForecastingSettings into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["forecastHorizon"]; ok {
		impl, err := unmarshalForecastHorizonImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ForecastHorizon' for 'ForecastingSettings': %+v", err)
		}
		s.ForecastHorizon = impl
	}

	if v, ok := temp["seasonality"]; ok {
		impl, err := unmarshalSeasonalityImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Seasonality' for 'ForecastingSettings': %+v", err)
		}
		s.Seasonality = impl
	}

	if v, ok := temp["targetLags"]; ok {
		impl, err := unmarshalTargetLagsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TargetLags' for 'ForecastingSettings': %+v", err)
		}
		s.TargetLags = impl
	}

	if v, ok := temp["targetRollingWindowSize"]; ok {
		impl, err := unmarshalTargetRollingWindowSizeImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TargetRollingWindowSize' for 'ForecastingSettings': %+v", err)
		}
		s.TargetRollingWindowSize = impl
	}
	return nil
}
