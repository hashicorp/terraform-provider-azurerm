package alertruletemplates

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AlertRuleTemplate = NrtAlertRuleTemplate{}

type NrtAlertRuleTemplate struct {
	Properties *NrtAlertRuleTemplateProperties `json:"properties,omitempty"`

	// Fields inherited from AlertRuleTemplate

	Id         *string                `json:"id,omitempty"`
	Kind       AlertRuleKind          `json:"kind"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

func (s NrtAlertRuleTemplate) AlertRuleTemplate() BaseAlertRuleTemplateImpl {
	return BaseAlertRuleTemplateImpl{
		Id:         s.Id,
		Kind:       s.Kind,
		Name:       s.Name,
		SystemData: s.SystemData,
		Type:       s.Type,
	}
}

var _ json.Marshaler = NrtAlertRuleTemplate{}

func (s NrtAlertRuleTemplate) MarshalJSON() ([]byte, error) {
	type wrapper NrtAlertRuleTemplate
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NrtAlertRuleTemplate: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NrtAlertRuleTemplate: %+v", err)
	}

	decoded["kind"] = "NRT"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NrtAlertRuleTemplate: %+v", err)
	}

	return encoded, nil
}
