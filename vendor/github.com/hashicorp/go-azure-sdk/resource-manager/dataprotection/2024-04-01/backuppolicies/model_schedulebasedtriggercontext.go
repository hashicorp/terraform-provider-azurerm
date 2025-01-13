package backuppolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TriggerContext = ScheduleBasedTriggerContext{}

type ScheduleBasedTriggerContext struct {
	Schedule        BackupSchedule    `json:"schedule"`
	TaggingCriteria []TaggingCriteria `json:"taggingCriteria"`

	// Fields inherited from TriggerContext

	ObjectType string `json:"objectType"`
}

func (s ScheduleBasedTriggerContext) TriggerContext() BaseTriggerContextImpl {
	return BaseTriggerContextImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = ScheduleBasedTriggerContext{}

func (s ScheduleBasedTriggerContext) MarshalJSON() ([]byte, error) {
	type wrapper ScheduleBasedTriggerContext
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ScheduleBasedTriggerContext: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ScheduleBasedTriggerContext: %+v", err)
	}

	decoded["objectType"] = "ScheduleBasedTriggerContext"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ScheduleBasedTriggerContext: %+v", err)
	}

	return encoded, nil
}
