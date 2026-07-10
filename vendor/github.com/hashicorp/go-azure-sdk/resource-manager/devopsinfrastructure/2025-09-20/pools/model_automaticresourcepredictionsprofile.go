package pools

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ResourcePredictionsProfile = AutomaticResourcePredictionsProfile{}

type AutomaticResourcePredictionsProfile struct {
	PredictionPreference *PredictionPreference `json:"predictionPreference,omitempty"`

	// Fields inherited from ResourcePredictionsProfile

	Kind ResourcePredictionsProfileType `json:"kind"`
}

func (s AutomaticResourcePredictionsProfile) ResourcePredictionsProfile() BaseResourcePredictionsProfileImpl {
	return BaseResourcePredictionsProfileImpl{
		Kind: s.Kind,
	}
}

var _ json.Marshaler = AutomaticResourcePredictionsProfile{}

func (s AutomaticResourcePredictionsProfile) MarshalJSON() ([]byte, error) {
	type wrapper AutomaticResourcePredictionsProfile
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AutomaticResourcePredictionsProfile: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AutomaticResourcePredictionsProfile: %+v", err)
	}

	decoded["kind"] = "Automatic"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AutomaticResourcePredictionsProfile: %+v", err)
	}

	return encoded, nil
}
