package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ VectorSearchCompressionConfiguration = ScalarQuantizationVectorSearchCompressionConfiguration{}

type ScalarQuantizationVectorSearchCompressionConfiguration struct {
	ScalarQuantizationParameters *ScalarQuantizationParameters `json:"scalarQuantizationParameters,omitempty"`

	// Fields inherited from VectorSearchCompressionConfiguration

	Kind                VectorSearchCompressionKind `json:"kind"`
	Name                string                      `json:"name"`
	RescoringOptions    *RescoringOptions           `json:"rescoringOptions,omitempty"`
	TruncationDimension *int64                      `json:"truncationDimension,omitempty"`
}

func (s ScalarQuantizationVectorSearchCompressionConfiguration) VectorSearchCompressionConfiguration() BaseVectorSearchCompressionConfigurationImpl {
	return BaseVectorSearchCompressionConfigurationImpl{
		Kind:                s.Kind,
		Name:                s.Name,
		RescoringOptions:    s.RescoringOptions,
		TruncationDimension: s.TruncationDimension,
	}
}

var _ json.Marshaler = ScalarQuantizationVectorSearchCompressionConfiguration{}

func (s ScalarQuantizationVectorSearchCompressionConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper ScalarQuantizationVectorSearchCompressionConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ScalarQuantizationVectorSearchCompressionConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ScalarQuantizationVectorSearchCompressionConfiguration: %+v", err)
	}

	decoded["kind"] = "scalarQuantization"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ScalarQuantizationVectorSearchCompressionConfiguration: %+v", err)
	}

	return encoded, nil
}
