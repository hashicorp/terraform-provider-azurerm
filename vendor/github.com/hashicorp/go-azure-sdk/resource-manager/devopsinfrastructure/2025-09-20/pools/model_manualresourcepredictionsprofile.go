package pools

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ResourcePredictionsProfile = ManualResourcePredictionsProfile{}

type ManualResourcePredictionsProfile struct {

	// Fields inherited from ResourcePredictionsProfile

	Kind ResourcePredictionsProfileType `json:"kind"`
}

func (s ManualResourcePredictionsProfile) ResourcePredictionsProfile() BaseResourcePredictionsProfileImpl {
	return BaseResourcePredictionsProfileImpl{
		Kind: s.Kind,
	}
}

var _ json.Marshaler = ManualResourcePredictionsProfile{}

func (s ManualResourcePredictionsProfile) MarshalJSON() ([]byte, error) {
	type wrapper ManualResourcePredictionsProfile
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ManualResourcePredictionsProfile: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ManualResourcePredictionsProfile: %+v", err)
	}

	decoded["kind"] = "Manual"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ManualResourcePredictionsProfile: %+v", err)
	}

	return encoded, nil
}
