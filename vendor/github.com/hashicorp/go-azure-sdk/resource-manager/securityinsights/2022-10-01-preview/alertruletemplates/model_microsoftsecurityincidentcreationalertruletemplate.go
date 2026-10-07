package alertruletemplates

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AlertRuleTemplate = MicrosoftSecurityIncidentCreationAlertRuleTemplate{}

type MicrosoftSecurityIncidentCreationAlertRuleTemplate struct {
	Properties *MicrosoftSecurityIncidentCreationAlertRuleTemplateProperties `json:"properties,omitempty"`

	// Fields inherited from AlertRuleTemplate

	Id         *string                `json:"id,omitempty"`
	Kind       AlertRuleKind          `json:"kind"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

func (s MicrosoftSecurityIncidentCreationAlertRuleTemplate) AlertRuleTemplate() BaseAlertRuleTemplateImpl {
	return BaseAlertRuleTemplateImpl{
		Id:         s.Id,
		Kind:       s.Kind,
		Name:       s.Name,
		SystemData: s.SystemData,
		Type:       s.Type,
	}
}

var _ json.Marshaler = MicrosoftSecurityIncidentCreationAlertRuleTemplate{}

func (s MicrosoftSecurityIncidentCreationAlertRuleTemplate) MarshalJSON() ([]byte, error) {
	type wrapper MicrosoftSecurityIncidentCreationAlertRuleTemplate
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MicrosoftSecurityIncidentCreationAlertRuleTemplate: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MicrosoftSecurityIncidentCreationAlertRuleTemplate: %+v", err)
	}

	decoded["kind"] = "MicrosoftSecurityIncidentCreation"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MicrosoftSecurityIncidentCreationAlertRuleTemplate: %+v", err)
	}

	return encoded, nil
}
