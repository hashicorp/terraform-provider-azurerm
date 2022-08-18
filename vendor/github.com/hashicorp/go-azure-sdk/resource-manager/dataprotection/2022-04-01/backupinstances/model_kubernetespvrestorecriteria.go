package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ItemLevelRestoreCriteria = KubernetesPVRestoreCriteria{}

type KubernetesPVRestoreCriteria struct {
	Name             *string `json:"name,omitempty"`
	StorageClassName *string `json:"storageClassName,omitempty"`

	// Fields inherited from ItemLevelRestoreCriteria
}

var _ json.Marshaler = KubernetesPVRestoreCriteria{}

func (s KubernetesPVRestoreCriteria) MarshalJSON() ([]byte, error) {
	type wrapper KubernetesPVRestoreCriteria
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KubernetesPVRestoreCriteria: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KubernetesPVRestoreCriteria: %+v", err)
	}
	decoded["objectType"] = "KubernetesPVRestoreCriteria"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KubernetesPVRestoreCriteria: %+v", err)
	}

	return encoded, nil
}
