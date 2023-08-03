package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = IsNullOrUndefinedAdvancedFilter{}

type IsNullOrUndefinedAdvancedFilter struct {

	// Fields inherited from AdvancedFilter
	Key *string `json:"key,omitempty"`
}

var _ json.Marshaler = IsNullOrUndefinedAdvancedFilter{}

func (s IsNullOrUndefinedAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper IsNullOrUndefinedAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling IsNullOrUndefinedAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling IsNullOrUndefinedAdvancedFilter: %+v", err)
	}
	decoded["operatorType"] = "IsNullOrUndefined"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling IsNullOrUndefinedAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
