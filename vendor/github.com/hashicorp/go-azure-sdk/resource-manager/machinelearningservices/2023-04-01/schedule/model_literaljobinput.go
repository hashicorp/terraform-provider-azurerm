package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobInput = LiteralJobInput{}

type LiteralJobInput struct {
	Value string `json:"value"`

	// Fields inherited from JobInput
	Description *string `json:"description,omitempty"`
}

var _ json.Marshaler = LiteralJobInput{}

func (s LiteralJobInput) MarshalJSON() ([]byte, error) {
	type wrapper LiteralJobInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LiteralJobInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LiteralJobInput: %+v", err)
	}
	decoded["jobInputType"] = "literal"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LiteralJobInput: %+v", err)
	}

	return encoded, nil
}
