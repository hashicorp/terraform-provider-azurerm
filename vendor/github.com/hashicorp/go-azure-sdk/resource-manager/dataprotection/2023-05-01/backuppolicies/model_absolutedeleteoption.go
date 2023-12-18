package backuppolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeleteOption = AbsoluteDeleteOption{}

type AbsoluteDeleteOption struct {

	// Fields inherited from DeleteOption
	Duration string `json:"duration"`
}

var _ json.Marshaler = AbsoluteDeleteOption{}

func (s AbsoluteDeleteOption) MarshalJSON() ([]byte, error) {
	type wrapper AbsoluteDeleteOption
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AbsoluteDeleteOption: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AbsoluteDeleteOption: %+v", err)
	}
	decoded["objectType"] = "AbsoluteDeleteOption"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AbsoluteDeleteOption: %+v", err)
	}

	return encoded, nil
}
