package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ScoringFunction = DistanceScoringFunction{}

type DistanceScoringFunction struct {
	Distance DistanceScoringParameters `json:"distance"`

	// Fields inherited from ScoringFunction

	Boost         float64                       `json:"boost"`
	FieldName     string                        `json:"fieldName"`
	Interpolation *ScoringFunctionInterpolation `json:"interpolation,omitempty"`
	Type          string                        `json:"type"`
}

func (s DistanceScoringFunction) ScoringFunction() BaseScoringFunctionImpl {
	return BaseScoringFunctionImpl{
		Boost:         s.Boost,
		FieldName:     s.FieldName,
		Interpolation: s.Interpolation,
		Type:          s.Type,
	}
}

var _ json.Marshaler = DistanceScoringFunction{}

func (s DistanceScoringFunction) MarshalJSON() ([]byte, error) {
	type wrapper DistanceScoringFunction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DistanceScoringFunction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DistanceScoringFunction: %+v", err)
	}

	decoded["type"] = "distance"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DistanceScoringFunction: %+v", err)
	}

	return encoded, nil
}
