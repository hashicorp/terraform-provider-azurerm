package backupinstanceresources

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ItemLevelRestoreCriteria = GenericRestoreDatasourceCriteria{}

type GenericRestoreDatasourceCriteria struct {
	ResourceSelectors ResourceListSelectionCriteria `json:"resourceSelectors"`

	// Fields inherited from ItemLevelRestoreCriteria

	ObjectType string `json:"objectType"`
}

func (s GenericRestoreDatasourceCriteria) ItemLevelRestoreCriteria() BaseItemLevelRestoreCriteriaImpl {
	return BaseItemLevelRestoreCriteriaImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = GenericRestoreDatasourceCriteria{}

func (s GenericRestoreDatasourceCriteria) MarshalJSON() ([]byte, error) {
	type wrapper GenericRestoreDatasourceCriteria
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GenericRestoreDatasourceCriteria: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GenericRestoreDatasourceCriteria: %+v", err)
	}

	decoded["objectType"] = "GenericRestoreDatasourceCriteria"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GenericRestoreDatasourceCriteria: %+v", err)
	}

	return encoded, nil
}
