package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CopyActivityTypeProperties struct {
	DataIntegrationUnits            *int64                           `json:"dataIntegrationUnits,omitempty"`
	EnableSkipIncompatibleRow       *bool                            `json:"enableSkipIncompatibleRow,omitempty"`
	EnableStaging                   *bool                            `json:"enableStaging,omitempty"`
	LogSettings                     *LogSettings                     `json:"logSettings,omitempty"`
	LogStorageSettings              *LogStorageSettings              `json:"logStorageSettings,omitempty"`
	ParallelCopies                  *int64                           `json:"parallelCopies,omitempty"`
	Preserve                        *[]string                        `json:"preserve,omitempty"`
	PreserveRules                   *[]string                        `json:"preserveRules,omitempty"`
	RedirectIncompatibleRowSettings *RedirectIncompatibleRowSettings `json:"redirectIncompatibleRowSettings,omitempty"`
	Sink                            CopySink                         `json:"sink"`
	SkipErrorFile                   *SkipErrorFile                   `json:"skipErrorFile,omitempty"`
	Source                          CopySource                       `json:"source"`
	StagingSettings                 *StagingSettings                 `json:"stagingSettings,omitempty"`
	Translator                      *interface{}                     `json:"translator,omitempty"`
	ValidateDataConsistency         *bool                            `json:"validateDataConsistency,omitempty"`
}

var _ json.Unmarshaler = &CopyActivityTypeProperties{}

func (s *CopyActivityTypeProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		DataIntegrationUnits            *int64                           `json:"dataIntegrationUnits,omitempty"`
		EnableSkipIncompatibleRow       *bool                            `json:"enableSkipIncompatibleRow,omitempty"`
		EnableStaging                   *bool                            `json:"enableStaging,omitempty"`
		LogSettings                     *LogSettings                     `json:"logSettings,omitempty"`
		LogStorageSettings              *LogStorageSettings              `json:"logStorageSettings,omitempty"`
		ParallelCopies                  *int64                           `json:"parallelCopies,omitempty"`
		Preserve                        *[]string                        `json:"preserve,omitempty"`
		PreserveRules                   *[]string                        `json:"preserveRules,omitempty"`
		RedirectIncompatibleRowSettings *RedirectIncompatibleRowSettings `json:"redirectIncompatibleRowSettings,omitempty"`
		SkipErrorFile                   *SkipErrorFile                   `json:"skipErrorFile,omitempty"`
		StagingSettings                 *StagingSettings                 `json:"stagingSettings,omitempty"`
		Translator                      *interface{}                     `json:"translator,omitempty"`
		ValidateDataConsistency         *bool                            `json:"validateDataConsistency,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.DataIntegrationUnits = decoded.DataIntegrationUnits
	s.EnableSkipIncompatibleRow = decoded.EnableSkipIncompatibleRow
	s.EnableStaging = decoded.EnableStaging
	s.LogSettings = decoded.LogSettings
	s.LogStorageSettings = decoded.LogStorageSettings
	s.ParallelCopies = decoded.ParallelCopies
	s.Preserve = decoded.Preserve
	s.PreserveRules = decoded.PreserveRules
	s.RedirectIncompatibleRowSettings = decoded.RedirectIncompatibleRowSettings
	s.SkipErrorFile = decoded.SkipErrorFile
	s.StagingSettings = decoded.StagingSettings
	s.Translator = decoded.Translator
	s.ValidateDataConsistency = decoded.ValidateDataConsistency

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling CopyActivityTypeProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["sink"]; ok {
		impl, err := UnmarshalCopySinkImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Sink' for 'CopyActivityTypeProperties': %+v", err)
		}
		s.Sink = impl
	}

	if v, ok := temp["source"]; ok {
		impl, err := UnmarshalCopySourceImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Source' for 'CopyActivityTypeProperties': %+v", err)
		}
		s.Source = impl
	}

	return nil
}
