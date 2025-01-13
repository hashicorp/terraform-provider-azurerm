package alertrules

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AlertRule = NrtAlertRule{}

type NrtAlertRule struct {
	Properties *NrtAlertRuleProperties `json:"properties,omitempty"`

	// Fields inherited from AlertRule

	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Kind       AlertRuleKind          `json:"kind"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

func (s NrtAlertRule) AlertRule() BaseAlertRuleImpl {
	return BaseAlertRuleImpl{
		Etag:       s.Etag,
		Id:         s.Id,
		Kind:       s.Kind,
		Name:       s.Name,
		SystemData: s.SystemData,
		Type:       s.Type,
	}
}

var _ json.Marshaler = NrtAlertRule{}

func (s NrtAlertRule) MarshalJSON() ([]byte, error) {
	type wrapper NrtAlertRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NrtAlertRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NrtAlertRule: %+v", err)
	}

	decoded["kind"] = "NRT"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NrtAlertRule: %+v", err)
	}

	return encoded, nil
}
