package alertrules

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AlertRule = FusionAlertRule{}

type FusionAlertRule struct {
	Properties *FusionAlertRuleProperties `json:"properties,omitempty"`

	// Fields inherited from AlertRule

	Etag       *string                `json:"etag,omitempty"`
	Id         *string                `json:"id,omitempty"`
	Kind       AlertRuleKind          `json:"kind"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

func (s FusionAlertRule) AlertRule() BaseAlertRuleImpl {
	return BaseAlertRuleImpl{
		Etag:       s.Etag,
		Id:         s.Id,
		Kind:       s.Kind,
		Name:       s.Name,
		SystemData: s.SystemData,
		Type:       s.Type,
	}
}

var _ json.Marshaler = FusionAlertRule{}

func (s FusionAlertRule) MarshalJSON() ([]byte, error) {
	type wrapper FusionAlertRule
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FusionAlertRule: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FusionAlertRule: %+v", err)
	}

	decoded["kind"] = "Fusion"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FusionAlertRule: %+v", err)
	}

	return encoded, nil
}
