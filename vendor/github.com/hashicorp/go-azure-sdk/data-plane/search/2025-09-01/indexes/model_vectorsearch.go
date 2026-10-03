package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VectorSearch struct {
	Algorithms   *[]VectorSearchAlgorithmConfiguration   `json:"algorithms,omitempty"`
	Compressions *[]VectorSearchCompressionConfiguration `json:"compressions,omitempty"`
	Profiles     *[]VectorSearchProfile                  `json:"profiles,omitempty"`
	Vectorizers  *[]VectorSearchVectorizer               `json:"vectorizers,omitempty"`
}

var _ json.Unmarshaler = &VectorSearch{}

func (s *VectorSearch) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Profiles *[]VectorSearchProfile `json:"profiles,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Profiles = decoded.Profiles

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling VectorSearch into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["algorithms"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Algorithms into list []json.RawMessage: %+v", err)
		}

		output := make([]VectorSearchAlgorithmConfiguration, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalVectorSearchAlgorithmConfigurationImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Algorithms' for 'VectorSearch': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Algorithms = &output
	}

	if v, ok := temp["compressions"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Compressions into list []json.RawMessage: %+v", err)
		}

		output := make([]VectorSearchCompressionConfiguration, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalVectorSearchCompressionConfigurationImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Compressions' for 'VectorSearch': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Compressions = &output
	}

	if v, ok := temp["vectorizers"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Vectorizers into list []json.RawMessage: %+v", err)
		}

		output := make([]VectorSearchVectorizer, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalVectorSearchVectorizerImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Vectorizers' for 'VectorSearch': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Vectorizers = &output
	}

	return nil
}
