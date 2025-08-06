package backupinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ItemLevelRestoreCriteria = ItemPathBasedRestoreCriteria{}

type ItemPathBasedRestoreCriteria struct {
	IsPathRelativeToBackupItem bool      `json:"isPathRelativeToBackupItem"`
	ItemPath                   string    `json:"itemPath"`
	SubItemPathPrefix          *[]string `json:"subItemPathPrefix,omitempty"`

	// Fields inherited from ItemLevelRestoreCriteria

	ObjectType string `json:"objectType"`
}

func (s ItemPathBasedRestoreCriteria) ItemLevelRestoreCriteria() BaseItemLevelRestoreCriteriaImpl {
	return BaseItemLevelRestoreCriteriaImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = ItemPathBasedRestoreCriteria{}

func (s ItemPathBasedRestoreCriteria) MarshalJSON() ([]byte, error) {
	type wrapper ItemPathBasedRestoreCriteria
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ItemPathBasedRestoreCriteria: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ItemPathBasedRestoreCriteria: %+v", err)
	}

	decoded["objectType"] = "ItemPathBasedRestoreCriteria"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ItemPathBasedRestoreCriteria: %+v", err)
	}

	return encoded, nil
}
