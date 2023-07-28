package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TriggerBase = CronTrigger{}

type CronTrigger struct {
	Expression string `json:"expression"`

	// Fields inherited from TriggerBase
	EndTime   *string `json:"endTime,omitempty"`
	StartTime *string `json:"startTime,omitempty"`
	TimeZone  *string `json:"timeZone,omitempty"`
}

var _ json.Marshaler = CronTrigger{}

func (s CronTrigger) MarshalJSON() ([]byte, error) {
	type wrapper CronTrigger
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CronTrigger: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CronTrigger: %+v", err)
	}
	decoded["triggerType"] = "Cron"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CronTrigger: %+v", err)
	}

	return encoded, nil
}
