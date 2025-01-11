package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ItemLevelRestoreCriteria = KubernetesStorageClassRestoreCriteria{}

type KubernetesStorageClassRestoreCriteria struct {
	Provisioner              *string `json:"provisioner,omitempty"`
	SelectedStorageClassName *string `json:"selectedStorageClassName,omitempty"`

	// Fields inherited from ItemLevelRestoreCriteria

	ObjectType string `json:"objectType"`
}

func (s KubernetesStorageClassRestoreCriteria) ItemLevelRestoreCriteria() BaseItemLevelRestoreCriteriaImpl {
	return BaseItemLevelRestoreCriteriaImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = KubernetesStorageClassRestoreCriteria{}

func (s KubernetesStorageClassRestoreCriteria) MarshalJSON() ([]byte, error) {
	type wrapper KubernetesStorageClassRestoreCriteria
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KubernetesStorageClassRestoreCriteria: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KubernetesStorageClassRestoreCriteria: %+v", err)
	}

	decoded["objectType"] = "KubernetesStorageClassRestoreCriteria"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KubernetesStorageClassRestoreCriteria: %+v", err)
	}

	return encoded, nil
}
