package pipelines

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Activity = CustomActivity{}

type CustomActivity struct {
	LinkedServiceName *LinkedServiceReference      `json:"linkedServiceName,omitempty"`
	Policy            *ActivityPolicy              `json:"policy,omitempty"`
	TypeProperties    CustomActivityTypeProperties `json:"typeProperties"`

	// Fields inherited from Activity

	DependsOn        *[]ActivityDependency     `json:"dependsOn,omitempty"`
	Description      *string                   `json:"description,omitempty"`
	Name             string                    `json:"name"`
	OnInactiveMarkAs *ActivityOnInactiveMarkAs `json:"onInactiveMarkAs,omitempty"`
	State            *ActivityState            `json:"state,omitempty"`
	Type             string                    `json:"type"`
	UserProperties   *[]UserProperty           `json:"userProperties,omitempty"`
}

func (s CustomActivity) Activity() BaseActivityImpl {
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

var _ json.Marshaler = CustomActivity{}

func (s CustomActivity) MarshalJSON() ([]byte, error) {
	type wrapper CustomActivity
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomActivity: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomActivity: %+v", err)
	}

	decoded["type"] = "Custom"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomActivity: %+v", err)
	}

	return encoded, nil
}
