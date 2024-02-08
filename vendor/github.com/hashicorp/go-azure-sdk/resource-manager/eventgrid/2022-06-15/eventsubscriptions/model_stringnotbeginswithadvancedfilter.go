package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AdvancedFilter = StringNotBeginsWithAdvancedFilter{}

type StringNotBeginsWithAdvancedFilter struct {
	Values *[]string `json:"values,omitempty"`

	// Fields inherited from AdvancedFilter
	Key *string `json:"key,omitempty"`
}

var _ json.Marshaler = StringNotBeginsWithAdvancedFilter{}

func (s StringNotBeginsWithAdvancedFilter) MarshalJSON() ([]byte, error) {
	type wrapper StringNotBeginsWithAdvancedFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StringNotBeginsWithAdvancedFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StringNotBeginsWithAdvancedFilter: %+v", err)
	}
	decoded["operatorType"] = "StringNotBeginsWith"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StringNotBeginsWithAdvancedFilter: %+v", err)
	}

	return encoded, nil
}
