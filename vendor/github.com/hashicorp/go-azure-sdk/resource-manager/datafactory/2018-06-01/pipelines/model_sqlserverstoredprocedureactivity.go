package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Activity = SqlServerStoredProcedureActivity{}

type SqlServerStoredProcedureActivity struct {
	LinkedServiceName *LinkedServiceReference                        `json:"linkedServiceName,omitempty"`
	Policy            *ActivityPolicy                                `json:"policy,omitempty"`
	TypeProperties    SqlServerStoredProcedureActivityTypeProperties `json:"typeProperties"`

	// Fields inherited from Activity

	DependsOn        *[]ActivityDependency     `json:"dependsOn,omitempty"`
	Description      *string                   `json:"description,omitempty"`
	Name             string                    `json:"name"`
	OnInactiveMarkAs *ActivityOnInactiveMarkAs `json:"onInactiveMarkAs,omitempty"`
	State            *ActivityState            `json:"state,omitempty"`
	Type             string                    `json:"type"`
	UserProperties   *[]UserProperty           `json:"userProperties,omitempty"`
}

func (s SqlServerStoredProcedureActivity) Activity() BaseActivityImpl {
	return BaseActivityImpl{
		DependsOn:        s.DependsOn,
		Description:      s.Description,
		Name:             s.Name,
		OnInactiveMarkAs: s.OnInactiveMarkAs,
		State:            s.State,
		Type:             s.Type,
		UserProperties:   s.UserProperties,
	}
}

var _ json.Marshaler = SqlServerStoredProcedureActivity{}

func (s SqlServerStoredProcedureActivity) MarshalJSON() ([]byte, error) {
	type wrapper SqlServerStoredProcedureActivity
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SqlServerStoredProcedureActivity: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SqlServerStoredProcedureActivity: %+v", err)
	}

	decoded["type"] = "SqlServerStoredProcedure"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SqlServerStoredProcedureActivity: %+v", err)
	}

	return encoded, nil
}
