package encodings

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobInput = JobInputSequence{}

type JobInputSequence struct {
	Inputs *[]JobInputClip `json:"inputs,omitempty"`

	// Fields inherited from JobInput
}

var _ json.Marshaler = JobInputSequence{}

func (s JobInputSequence) MarshalJSON() ([]byte, error) {
	type wrapper JobInputSequence
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JobInputSequence: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JobInputSequence: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.JobInputSequence"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JobInputSequence: %+v", err)
	}

	return encoded, nil
}
