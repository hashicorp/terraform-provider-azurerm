package synchronizationsetting

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SynchronizationSetting = ScheduledSynchronizationSetting{}

type ScheduledSynchronizationSetting struct {
	Properties ScheduledSynchronizationSettingProperties `json:"properties"`

	// Fields inherited from SynchronizationSetting
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`
}

var _ json.Marshaler = ScheduledSynchronizationSetting{}

func (s ScheduledSynchronizationSetting) MarshalJSON() ([]byte, error) {
	type wrapper ScheduledSynchronizationSetting
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ScheduledSynchronizationSetting: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ScheduledSynchronizationSetting: %+v", err)
	}
	decoded["kind"] = "ScheduleBased"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ScheduledSynchronizationSetting: %+v", err)
	}

	return encoded, nil
}
