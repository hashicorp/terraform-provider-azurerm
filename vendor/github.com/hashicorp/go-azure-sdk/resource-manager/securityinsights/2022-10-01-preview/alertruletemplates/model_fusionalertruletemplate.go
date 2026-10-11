package alertruletemplates

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AlertRuleTemplate = FusionAlertRuleTemplate{}

type FusionAlertRuleTemplate struct {
	Properties *FusionAlertRuleTemplateProperties `json:"properties,omitempty"`

	// Fields inherited from AlertRuleTemplate

	Id         *string                `json:"id,omitempty"`
	Kind       AlertRuleKind          `json:"kind"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

func (s FusionAlertRuleTemplate) AlertRuleTemplate() BaseAlertRuleTemplateImpl {
	return BaseAlertRuleTemplateImpl{
		Id:         s.Id,
		Kind:       s.Kind,
		Name:       s.Name,
		SystemData: s.SystemData,
		Type:       s.Type,
	}
}

var _ json.Marshaler = FusionAlertRuleTemplate{}

func (s FusionAlertRuleTemplate) MarshalJSON() ([]byte, error) {
	type wrapper FusionAlertRuleTemplate
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FusionAlertRuleTemplate: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FusionAlertRuleTemplate: %+v", err)
	}

	decoded["kind"] = "Fusion"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FusionAlertRuleTemplate: %+v", err)
	}

	return encoded, nil
}
