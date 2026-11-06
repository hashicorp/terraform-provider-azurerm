package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ScoringFunction = MagnitudeScoringFunction{}

type MagnitudeScoringFunction struct {
	Magnitude MagnitudeScoringParameters `json:"magnitude"`

	// Fields inherited from ScoringFunction

	Boost         float64                       `json:"boost"`
	FieldName     string                        `json:"fieldName"`
	Interpolation *ScoringFunctionInterpolation `json:"interpolation,omitempty"`
	Type          string                        `json:"type"`
}

func (s MagnitudeScoringFunction) ScoringFunction() BaseScoringFunctionImpl {
	return BaseScoringFunctionImpl{
		Boost:         s.Boost,
		FieldName:     s.FieldName,
		Interpolation: s.Interpolation,
		Type:          s.Type,
	}
}

var _ json.Marshaler = MagnitudeScoringFunction{}

func (s MagnitudeScoringFunction) MarshalJSON() ([]byte, error) {
	type wrapper MagnitudeScoringFunction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MagnitudeScoringFunction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MagnitudeScoringFunction: %+v", err)
	}

	decoded["type"] = "magnitude"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MagnitudeScoringFunction: %+v", err)
	}

	return encoded, nil
}
