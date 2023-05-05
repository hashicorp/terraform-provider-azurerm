package protectionpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SchedulePolicy = LongTermSchedulePolicy{}

type LongTermSchedulePolicy struct {

	// Fields inherited from SchedulePolicy
}

var _ json.Marshaler = LongTermSchedulePolicy{}

func (s LongTermSchedulePolicy) MarshalJSON() ([]byte, error) {
	type wrapper LongTermSchedulePolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LongTermSchedulePolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LongTermSchedulePolicy: %+v", err)
	}
	decoded["schedulePolicyType"] = "LongTermSchedulePolicy"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LongTermSchedulePolicy: %+v", err)
	}

	return encoded, nil
}
