package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobLimits = CommandJobLimits{}

type CommandJobLimits struct {

	// Fields inherited from JobLimits
	Timeout *string `json:"timeout,omitempty"`
}

var _ json.Marshaler = CommandJobLimits{}

func (s CommandJobLimits) MarshalJSON() ([]byte, error) {
	type wrapper CommandJobLimits
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CommandJobLimits: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CommandJobLimits: %+v", err)
	}
	decoded["jobLimitsType"] = "Command"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CommandJobLimits: %+v", err)
	}

	return encoded, nil
}
