package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ VectorSearchCompressionConfiguration = BinaryQuantizationVectorSearchCompressionConfiguration{}

type BinaryQuantizationVectorSearchCompressionConfiguration struct {

	// Fields inherited from VectorSearchCompressionConfiguration

	Kind                VectorSearchCompressionKind `json:"kind"`
	Name                string                      `json:"name"`
	RescoringOptions    *RescoringOptions           `json:"rescoringOptions,omitempty"`
	TruncationDimension *int64                      `json:"truncationDimension,omitempty"`
}

func (s BinaryQuantizationVectorSearchCompressionConfiguration) VectorSearchCompressionConfiguration() BaseVectorSearchCompressionConfigurationImpl {
	return BaseVectorSearchCompressionConfigurationImpl{
		Kind:                s.Kind,
		Name:                s.Name,
		RescoringOptions:    s.RescoringOptions,
		TruncationDimension: s.TruncationDimension,
	}
}

var _ json.Marshaler = BinaryQuantizationVectorSearchCompressionConfiguration{}

func (s BinaryQuantizationVectorSearchCompressionConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper BinaryQuantizationVectorSearchCompressionConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BinaryQuantizationVectorSearchCompressionConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BinaryQuantizationVectorSearchCompressionConfiguration: %+v", err)
	}

	decoded["kind"] = "binaryQuantization"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BinaryQuantizationVectorSearchCompressionConfiguration: %+v", err)
	}

	return encoded, nil
}
