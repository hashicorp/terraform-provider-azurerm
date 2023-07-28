package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ NCrossValidations = CustomNCrossValidations{}

type CustomNCrossValidations struct {
	Value int64 `json:"value"`

	// Fields inherited from NCrossValidations
}

var _ json.Marshaler = CustomNCrossValidations{}

func (s CustomNCrossValidations) MarshalJSON() ([]byte, error) {
	type wrapper CustomNCrossValidations
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomNCrossValidations: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomNCrossValidations: %+v", err)
	}
	decoded["mode"] = "Custom"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomNCrossValidations: %+v", err)
	}

	return encoded, nil
}
