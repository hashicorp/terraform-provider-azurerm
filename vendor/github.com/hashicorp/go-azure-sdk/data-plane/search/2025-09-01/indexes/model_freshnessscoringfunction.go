package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ScoringFunction = FreshnessScoringFunction{}

type FreshnessScoringFunction struct {
	Freshness FreshnessScoringParameters `json:"freshness"`

	// Fields inherited from ScoringFunction

	Boost         float64                       `json:"boost"`
	FieldName     string                        `json:"fieldName"`
	Interpolation *ScoringFunctionInterpolation `json:"interpolation,omitempty"`
	Type          string                        `json:"type"`
}

func (s FreshnessScoringFunction) ScoringFunction() BaseScoringFunctionImpl {
	return BaseScoringFunctionImpl{
		Boost:         s.Boost,
		FieldName:     s.FieldName,
		Interpolation: s.Interpolation,
		Type:          s.Type,
	}
}

var _ json.Marshaler = FreshnessScoringFunction{}

func (s FreshnessScoringFunction) MarshalJSON() ([]byte, error) {
	type wrapper FreshnessScoringFunction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FreshnessScoringFunction: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FreshnessScoringFunction: %+v", err)
	}

	decoded["type"] = "freshness"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FreshnessScoringFunction: %+v", err)
	}

	return encoded, nil
}
